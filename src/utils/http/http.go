// 一些小便利
package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/serenize/snaker"
)

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

// 将http请求解析到结构体中
func Unpack(req *http.Request, ptr interface{}) error {

	if err := req.ParseForm(); err != nil {
		return err
	}

	fields := make(map[string]reflect.Value)

	// 这里默认 ptr 是指针类型，通过Elem方法取到其指向的值
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		if name == "" {
			//			name = strings.ToLower(fieldInfo.Name)
			name = snaker.CamelToSnake(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	for name, values := range req.Form {
		f := fields[name]

		if !f.IsValid() {
			continue // 忽略不认识的参数
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice { // 如果是数组类型，则拼成数组
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
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

// 将一个string类型填到任意类型
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
