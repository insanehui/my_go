package json

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"utils"
)

type I interface{}
type Obj map[string]I
type jres []Obj

// 对转化json的一个简易封装
func ToJson(data interface{}) string {
	j, _ := json.MarshalIndent(data, "", "    ")
	return string(j)
}

// 将 a append 到 s 当中去，返回一个新的slice
func Append(s interface{}, a interface{}) []interface{} {
	v := reflect.ValueOf(s) // 取到反射值
	var ret []interface{}

	// 判断类型
	switch v.Kind() {
	case reflect.Slice: // 如果是数组
		for i := 0; i < v.Len(); i++ {
			ret = append(ret, v.Index(i).Interface())
		}
	case reflect.Map:
		log.Printf("not support Map")
		break
	default:
		if s != nil {
			ret = append(ret, s)
		}
	}
	ret = append(ret, a)
	return ret
}

// 将集合转成set
func ToSet(s []string) map[string]int {

	ret := make(map[string]int)
	for _, v := range s {
		ret[v] = 1
	}
	return ret
}

func SetHas(set map[string]int, key string) bool {
	_, ok := set[key]
	return ok
}

// 命名参考了lodash
func OmitNew(j map[string]interface{}, key_str string) Obj {
	// 注：这里的类型，不能直接用 Obj？
	// 那自定义类型的用处在哪？

	// 类似以前定义的 json_unfilter之类的
	keys := strings.Split(key_str, " ")

	// 将keys存到一个set（用map来实现）里面
	key_set := make(map[string]int)
	for _, key := range keys {
		key_set[key] = 1
	}

	ret := make(Obj)
	// 将msg的字段过滤到jmsg
	for k, v := range j {
		if _, ok := key_set[k]; !ok {
			ret[k] = v
		}
	}

	return ret
}

// 只能再写一个纯 interface 版
func Omit(j interface{}, bys []string) Obj {
	// j 形如: Obj类型

	ret := make(Obj)

	rv := reflect.ValueOf(j)
	switch rv.Kind() {
	case reflect.Map:
		by_set := ToSet(bys)

		// 遍历 j
		for _, rkey := range rv.MapKeys() {
			key := rkey.String()
			if !SetHas(by_set, key) {
				ret[key] = rv.MapIndex(rkey).Interface()
			}
		}
	}
	return ret
}

// json group
func JsonGroup(o jres, by []string, opt ...string) Obj {

	ret := make(Obj)
	is_nostrip := false

	for _, op := range opt {
		if op == "nostrip" {
			is_nostrip = true
		}
	}

	for _, item := range o {
		// 根据 by 中的key，重新构造一个树结构，用以“收留”原数据中的元素
		prev := ret
		for i, key := range by {
			val := utils.ToStr(item[key]) // 硬转为string类型的val
			if i < len(by)-1 {            // 如果还没到“叶子”，继续延伸“树枝”

				// 先判断该“树枝”是否已经存在
				if br, ok := prev[val]; ok {
					v := reflect.ValueOf(br)
					switch v.Kind() {
					case reflect.Map: // 这是期望的 map 的情况
						prev = v.Interface().(Obj)
					}

				} else {
					tmp := make(Obj)
					prev[val] = tmp
					prev = tmp
				}
			} else { // 已经到了“叶子”，直接赋值
				leaf := make(Obj)
				if is_nostrip {
					leaf = item
				} else {
					leaf = Omit(item, by)
				}
				prev[val] = Append(prev[val], leaf)
			}
			//			log.Printf("%+v, %+v: %+v", r, i, ToJson(ret))

		}
	}
	return ret

}

// JsonGroup的替代
func Group(o interface{}, by []string, opt ...string) Obj {

	ret := make(Obj)
	is_nostrip := false

	for _, op := range opt {
		if op == "nostrip" {
			is_nostrip = true
		}
	}

	rO := reflect.ValueOf(o)
	switch rO.Kind() {
	case reflect.Slice:
		for i := 0; i < rO.Len(); i++ {
			rI := rO.Index(i)

			prev := ret
			for i, key := range by {

				item := rI.Interface()
				val := utils.ToStr(rI.MapIndex(reflect.ValueOf(key))) // 硬转为string类型的val
				if i < len(by)-1 {                                    // 如果还没到“叶子”，继续延伸“树枝”

					// 先判断该“树枝”是否已经存在
					if br, ok := prev[val]; ok {
						v := reflect.ValueOf(br)
						switch v.Kind() {
						case reflect.Map: // 这是期望的 map 的情况
							prev = v.Interface().(Obj)
						}

					} else {
						tmp := make(Obj)
						prev[val] = tmp
						prev = tmp
					}
				} else { // 已经到了“叶子”，直接赋值
					var leaf interface{}
					if is_nostrip {
						leaf = item
					} else {
						leaf = Omit(item, by)
					}
					prev[val] = Append(prev[val], leaf)
				}
			}

		}
	default:
		log.Printf("o is not Slice")

	}
	return ret
}
