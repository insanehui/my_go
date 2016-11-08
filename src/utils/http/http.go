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
	// V "github.com/asaskevich/govalidator"
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
// 代码摘自《go语言圣经》 ch12/params
func Unpack(req *http.Request, ptr interface{}) error {

	// 将输入参数解析到 req.Form 对象里
	if err := req.ParseForm(); err != nil {
		return err
	}

	// 临时存储 struct 里的数据
	fields := make(map[string]reflect.Value)

	// 这里默认 ptr 是指针，通过Elem方法取到其指向的值
	v := reflect.ValueOf(ptr).Elem()

	// 默认其是一个struct，遍历其所有字段
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // 取到字段的一些元数据, reflect.StructField
		tag := fieldInfo.Tag           // 首先取到字段tag. reflect.StructTag
		name := tag.Get("http")		// 取到tag中http的属性
		if name == "" {
			//			name = strings.ToLower(fieldInfo.Name)
			name = snaker.CamelToSnake(fieldInfo.Name) // 转为snaker形式的name
		}
		fields[name] = v.Field(i) // 存到map里的是rv（reflect.Value）
	}

	// 遍历请求的参数
	for name, values := range req.Form {

		f := fields[name] // 取到当前参数对应的field字段

		if !f.IsValid() {
			continue 
		}

		for _, value := range values {
			if f.Kind() == reflect.Slice { // 如果是数组类型，则拼成数组

				elem := reflect.New(f.Type().Elem()).Elem() // 这里有点晦涩啊...

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
