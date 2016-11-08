package utils

import (
	"crypto/md5"
	"reflect"
	//	"strings"
	//"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	S "github.com/fatih/structs"
)

// 小写形式
func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func MD5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%X", hashMd5.Sum(nil))
}

// 转字符串
func ToStr(s interface{}) string {
	return fmt.Sprintf("%+v", s)
}

func Keys(m interface{}) []interface{} {
	var ret []interface{}

	rv := reflect.ValueOf(m)
	switch rv.Kind() {
	case reflect.Map:

		// 遍历
		for _, rkey := range rv.MapKeys() {
			key := rkey.Interface()
			ret = append(ret, key)
		}
	}

	return ret
}

// 类似Keys，但key限定为string
func KeyStrs(m interface{}) []string {
	var ret []string

	rv := reflect.ValueOf(m)
	switch rv.Kind() {
	case reflect.Map:

		// 遍历
		for _, rkey := range rv.MapKeys() {
			key := rkey.String()
			ret = append(ret, key)
		}
	case reflect.Struct:
		return S.Names(m)
	}

	return ret
}

// 通用的对象转换
func Conv(from interface{}, to interface{}) error {
	j, err := json.Marshal(from)
	if err != nil {
		return err
	}

	err = json.Unmarshal(j, to)
	if err != nil {
		return err
	}

	return nil
}

//// JSON系列
//type JSON map[string]interface{}

//// omit，命名参考lodash
//// TODO 后面考虑改写成不定参数
//func (m JSON) OmitNew(key_str string) JSON {
//	keys := strings.Split(key_str, " ")

//	// 将keys存到一个set（用map来实现）里面
//	key_set := make(map[string]int)
//	for _, key := range keys {
//		key_set[key] = 1
//	}

//	ret := make(JSON)
//	// 将msg的字段过滤到jmsg
//	for k, v := range m {
//		if _, ok := key_set[k]; !ok {
//			ret[k] = v
//		}
//	}

//	return ret
//}
