// 一些小便利
package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	rt "utils/runtime"
	U "utils"

	V "github.com/asaskevich/govalidator"
	"github.com/serenize/snaker"
)

func log_client(r *http.Request) {
	log.Printf("<<<<<<<<<<<<<<<< %+v <<<<<<<<<<<<<<<<<<<", r.RemoteAddr)
	// log.Printf("[%+v]", r.RequestURI)
	log.Printf("[%+v]", r.URL)
}

type IRetJson interface {
	FromPanic(interface{})
}

// para为自定义的传入参数
// ret为自定义的返回参数
// 使用示例:
// func test(w http.ResponseWriter, r *http.Request) {
// 	var para struct {
// 		// ...
// 	}
//	var ret struct {
//		// ...
//	}
//	JsonDo(w, r, &para, &ret, func(){
//		// 处理para和ret
//	})
//
// }
func JsonDo(w http.ResponseWriter, r *http.Request, para interface{}, ret IRetJson, fn func()) {

	log_client(r)

	defer func() {
		p := recover()
		rt.Log(p)
		ret.FromPanic(p)
		WriteJson(w, ret)
	}()

	Checkout_(r, para)
	fn()
}

// 与JsonDo类似，但更通用（不一定返回json）
// end(p) 将在 defer里执行, 参数p 为panic
// 使用示例:
// func test(w http.ResponseWriter, r *http.Request) {
// 	var para struct {
// 		// ...
// 	}
//	var ret struct {
//		// ...
//	}
//	Do(w, r, &para, &ret, func(){
//		// 处理para和ret
//	},func(p){
//		// 处理panic
//		// 返回http回应
//	})
//
// }
func Do(w http.ResponseWriter, r *http.Request, para interface{}, ret interface{}, fn func(), end func(p interface{})) {
	log_client(r)

	defer func() {
		p := recover()
		rt.Log(p)
		end(p)
	}()

	Checkout_(r, para)
	fn()
}

func init() {
	V.SetFieldsRequiredByDefault(true) // 这是validator常用的设置
}

// 设置http头的方法
func Json(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json;charset=utf-8")
}

func Plain(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
}

// 直接返回一个json
func WriteJson(w http.ResponseWriter, data interface{}) {
	Json(w)
	ret, _ := json.MarshalIndent(data, "", "    ")
	w.Write(ret)
}

// 将http请求解析到结构体中. 使用tag: http, def
// 代码摘自《go语言圣经》 ch12/params
func Unpack(req *http.Request, ptr interface{}) error {

	// 将输入参数解析到 req.Form 对象里
	if err := req.ParseForm(); err != nil {
		return err
	}

	log.Printf("[Client Form:]:\n%+v", req.Form)

	// 临时存储 struct 里的数据
	fields := make(map[string]reflect.Value)

	// 这里默认 ptr 是指针，通过Elem方法取到其指向的值
	v := reflect.ValueOf(ptr).Elem()

	// 默认其是一个struct，遍历其所有字段
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // 取到字段的一些元数据, reflect.StructField
		tag := fieldInfo.Tag           // 首先取到字段tag. reflect.StructTag
		name := tag.Get("http")        // 取到tag中http的属性
		if name == "" {
			//			name = strings.ToLower(fieldInfo.Name)
			name = snaker.CamelToSnake(fieldInfo.Name) // 转为snaker形式的name
		}

		f := v.Field(i)
		fields[name] = f // 存到map里的是rv（reflect.Value）

		// 先填缺省值
		def := tag.Get("def")
		if def != "" {
			populate(f, def)
		}
	}

	// 遍历请求的参数
	for name, values := range req.Form {

		// f为当前字段的rv
		f := fields[name]
		if !f.IsValid() {
			continue
		}

		for _, value := range values {
			if f.Kind() == reflect.Slice { // 如果是数组类型，则拼成数组

				// new一个其元素对应的rv类型
				elem := reflect.New(f.Type().Elem()).Elem()

				// 将value（string）填到elem
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}

				// reflect世界里的append
				f.Set(reflect.Append(f, elem))

			} else { // 否则，后面值将前面的覆盖
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}

	return nil
}

// 在Unpack的基础上，用govalidator来校验参数合法性
func Checkout(req *http.Request, ptr interface{}) error {

	err := Unpack(req, ptr)
	if err != nil {
		return err
	}

	// ptr里已经填好值，现用validator对其进行校验
	_, err = V.ValidateStruct(ptr)
	if err != nil {
		return fmt.Errorf("params invalid: %s", err.Error())
	}

	return nil
}

// 抛异常的版本
func Checkout_(req *http.Request, ptr interface{}) {

	err := Checkout(req, ptr)
	if err != nil {
		panic(err)
	}
}

// 将一个string类型填到任意类型, 给Unpack使用，同样摘自《Go语言圣经》
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Uint, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	default:
		log.Printf("不支持类型：%s", v.Type())
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

// Post, 返回string
func Post(site string, vals map[string]interface{}) (res string, err error) {
	// TODO: 后续支持将struct作为参数传入

	url_vals := toUrlVals(vals)

	resp, err := http.PostForm(site, url_vals)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	res = string(body)
	return
}

// map -> url.Values
func toUrlVals(m map[string]interface{}) url.Values {
	ret := make(url.Values)
	for k, v := range m {
		str := U.ToStr(v)
		ss := []string{str}
		ret[k] = ss
	}
	return ret
}

