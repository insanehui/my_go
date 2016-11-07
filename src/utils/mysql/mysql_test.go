package mysql

import (
	"log"
	"reflect"
	"testing"
	j "utils/json"
)

func db() *DB {
	d, _ := Open("mysql", "teamtalk:12345@tcp(115.29.233.2:3306)/test?charset=utf8")
	return d
}

func TestQuery(t *testing.T) {
	db, _ := Open("mysql", "teamtalk:12345@tcp(115.29.233.2:3306)/test?charset=utf8")

	{
		res := db.Query("select * from user where id > ?", 1)
		log.Printf("%s", j.ToJson(res))
	}
	{
		res := db.Query("select count(*) from user where id > ?", 1)
		log.Printf("%s", j.ToJson(res))
	}
}

func Test_parseKey(t *testing.T) {
	d := map[string]int{"a": 1, "b": 2, "v": 3}
	{
		strs, vals := parseKey(d, "a", "v")
		log.Printf("%+v", strs)
		log.Printf("%+v", vals)
	}
	{
		strs, vals := parseKey(d)
		log.Printf("%+v", strs)
		log.Printf("%+v", vals)
	}

}

func TestExist(t *testing.T) {
	db, _ := Open("mysql", "teamtalk:12345@tcp(115.29.233.2:3306)/test?charset=utf8")
	d := map[string]string{"name": "haha"}
	ret := db.Exist("user", d, "name")
	log.Printf("%s", j.ToJson(ret))
}

func TestCount(t *testing.T) {
	d := db()
	m := map[string]string{"name": "haha"}
	c := d.Count("user", m, "name")
	if c != 3 {
		t.Error("Count fail")
	}
}

func TestQryValue(t *testing.T) {

	d := db()
	v := d.QryValue("select count(*) from user")
	log.Printf("%+v", v)

}

func TestQryInt(t *testing.T) {

	d := db()
	v := d.QryInt("select count(*) from user")
	log.Printf("type: %s", reflect.ValueOf(v).Kind())

}

func Test_insert(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("recover: %+v", p)
		}
	}()

	d := db()
	{
		m := map[string]interface{}{"id": 9, "name": "jialiao", "intval": 777}
		r := d.insert("user", m, true)
		log.Printf("affected: %+v", r)
	}
	{
		m := map[string]interface{}{"id": 9, "name": "jialiao", "intval": 777}
		r := d.insert("user", m, false)
		log.Printf("affected: %+v", r)
	}

}

func TestUpdate(t *testing.T) {
	d := db()
	m := map[string]interface{}{"id": 1, "name": "new_name", "intval": 88}
	r := d.Update("user", m, "id")
	log.Printf("affected: %+v", r)
}

func TestSet(t *testing.T) {
	d := db()
	m := map[string]interface{}{"name": "xxxx", "intval": 100}
	r := d.Set("user", m, "name")
	log.Printf("affected: %+v", r)
}
