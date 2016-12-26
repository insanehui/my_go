package mysql

import (
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"

	U "utils"
	// S "github.com/fatih/structs"

	_ "github.com/go-sql-driver/mysql"
)

type I interface{}
type Row map[string]I
type Res []Row

// 感觉是工厂模式才这样开接口，这里不开工厂，故不需要定义接口类型

//type IDB interface {
//	//	Count(table string, keys []string) uint64
//	Query(stmt string, paras ...interface{}) Res
//}

type DB struct {
	d *sql.DB
}

func Open(a string, b string) (*DB, error) {
	db := new(DB)
	var err error
	db.d, err = sql.Open(a, b)
	return db, err
}

// 示例：
// 格式：用户名:密码@地址/默认数据库?其他设置
// db = Mysql.Open_("mysql", "blueprint:ctg123@tcp(10.10.12.2:3306)/blueprint?charset=utf8")
func Open_(a string, b string) *DB {

	db, err := Open(a, b)
	if err != nil {
		log.Printf("open err: %+v", err) // 事实上，对于网络之类的问题，不会在此处抛异常
		panic(err)
	}
	return db
}

// 通过sql语句查询
func (me *DB) Query(stmtStr string, paras ...interface{}) Res {
	var res Res
	db := me.d

	rows, err := db.Query(stmtStr, paras...)

	if err != nil {
		log.Printf("err: %+v", err)
	}

	columns, _ := rows.Columns()
	count := len(columns)
	// 注：这是go很麻烦的地方。值和指针，要分开定义
	// 为什么会出现如此不合常理的设计？
	// 这恰恰是go对泛型支持的缺乏造成的
	// mysql的 rows.Scan 只接受 interface{}
	// 但又必须指针类型才能进行 *取值 操作
	// 而go又不支持类型转换。因此常用的做法便是：
	// 每个类型都定义相应的变量“代表”来参与相关行为。
	// 故而出现了需要同时定义 values 和 valuePtrs 的现象
	// go致力于将代码变得简单优雅，但这样的设计恰恰与简单优雅背道而驰
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {

		// 建立指针 和 值 的对应关系
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		// scan只接受指针
		rows.Scan(valuePtrs...)
		row := make(Row)

		for i, col := range columns {
			row[col] = value(values[i])
		}

		res = append(res, row)
	}

	return res
}

// 查询单个值
func (me *DB) QryValue(q string, paras ...interface{}) interface{} {
	db := me.d
	row := db.QueryRow(q, paras...)
	var ret interface{}
	row.Scan(&ret)
	return value(ret)
}

// 查询int值
func (me *DB) QryInt(q string, paras ...interface{}) int64 {
	// TODO 姑且先用 int64，后面再看看如何处理 int32 的问题
	v := me.QryValue(q, paras...)
	rv := reflect.ValueOf(v)
	var ret int64
	switch rv.Kind() {
	case reflect.String:
		ret, _ = strconv.ParseInt(rv.String(), 10, 64)
	default:
		ret = rv.Int()
	}
	return ret
}

// 计数（与exist逻辑类型，但仅统计）
func (me *DB) Count(table string, j interface{}, keys ...interface{}) int64 {
	q := "select count(*) from " + table + " where "
	exprs, vals := parseKey(j, keys...)
	q += strings.Join(exprs, " and ")
	return me.QryInt(q, vals...)
}

func (me *DB) Exist(table string, j interface{}, keys ...interface{}) Res {
	// 为什么返回的是int，而不是uint64呢，出于以下考虑：
	// 如果count的值超过了count的范围，这整个数据库得重新设计了
	q := "select * from " + table + " where "
	exprs, vals := parseKey(j, keys...)
	q += strings.Join(exprs, " and ")
	ret := me.Query(q, vals...)

	return ret
}

func (me *DB) Insert(table string, data interface{}) int64 {
	return me.insert(table, data, false)
}

func (me *DB) Replace(table string, data interface{}) int64 {
	return me.insert(table, data, true)
}

func (me *DB) Update(table string, data interface{}, keys ...interface{}) int64 {
	q := "update " + table + " set "
	exprs, vals := parseKey(data)
	q += strings.Join(exprs, ",")
	q += " where "
	conExprs, conVals := parseKey(data, keys...)
	q += strings.Join(conExprs, " and ")

	paras := append(vals, conVals...)
	rows, _ := me.Exec(q, paras...)
	return rows

}

func (me *DB) Set(table string, data interface{}, keys ...interface{}) int64 {
	n := me.Count(table, data, keys...)
	if n == 0 {
		return me.Insert(table, data)
	} else {
		return me.Update(table, data, keys...)
	}
}

// ==============================================================================

// 内部insert函数，返回的是last insert_id
func (me *DB) insert(table string, data interface{}, replace bool) int64 {
	var op string
	if replace {
		op = "replace" // 如果replace生效，并产生修改，affected的值会为2
	} else {
		op = "insert"
	}
	q := op + " into " + table + " set "
	exprs, vals := parseKey(data)
	q += strings.Join(exprs, ",")

	_, last_id := me.Exec(q, vals...)
	return last_id
}

func (me *DB) Exec(q string, args ...interface{}) (int64, int64) {
	db := me.d
	var err error
	var res sql.Result
	log.Printf("sql: %s, args: %+v", q, args)
	res, err = db.Exec(q, args...)
	checkErr(err)
	rows, err := res.RowsAffected()
	checkErr(err)
	last_id, err := res.LastInsertId()
	return rows, last_id
}

func parseKey(_d interface{}, keys ...interface{}) ([]string, []interface{}) {
	var exprs []string
	var vals []interface{}

	var d interface{}

	// 先用反射来检查传入d的类型
	rv := reflect.ValueOf(_d)

	// 如果是struct，将其转成map
	if rv.Kind() == reflect.Struct {

		m := make(map[string]interface{})
		U.Conv(_d, &m)
		rv = reflect.ValueOf(m)
		d = m
	} else {
		d = _d
	}

	// 如果没有传keys，则取map的所有key
	if len(keys) == 0 {
		keys = U.Keys(d)
	}

	for _, key := range keys {
		exprs = append(exprs, U.ToStr(key)+" = ? ")

		vals = append(vals, rv.MapIndex(reflect.ValueOf(key)).Interface())
	}

	return exprs, vals
}

// 转换sql结果中的值（这就是一个大坑）
func value(a interface{}) interface{} {

	// 注：（更大的坑）：
	// go的mysql接口查出来的数据，类型还是不固定的！！
	var v interface{}

	b, ok := a.([]byte)

	if ok {
		v = string(b)
		//		log.Printf("col type: [string]")
	} else {
		v = a
		//		log.Printf("col type: [int?]")
	}
	return v
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

// 纯粹用来将“log”使用上而已
func test() {
	log.Printf("haha")
}
