package yaml

import (
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
	// "gopkg.in/yaml.v2"

	I "utils/io"
	// I "../io"
)

// 从文件解析一个yaml
func FromFile(filename string) interface{} {
	// res := make(map[string]interface{})
	var res interface{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("read file error: %+v", err.Error())
		return res
	}

	err = yaml.Unmarshal(data, &res)
	if err != nil {
		log.Printf("yaml parse error: %+v", err.Error())
		return res
	}

	return res
}

// 解析到一个目的变量（可以是一个interface{}，或者struct, slice等）
func FromFileTo_(filename string, to interface{}) {
	b := I.ReadFile_(filename)
	err := yaml.Unmarshal(b, to)
	if err != nil {
		log.Panicf("yaml.FromFileTo err: %+v", err.Error())
	}
}
