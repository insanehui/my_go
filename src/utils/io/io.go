package io

import (
	"io/ioutil"
	"log"
)

func ReadFile_(filename string) []byte {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Panicf("read file err: %+v", err.Error())
	}

	return data
}

// 写小文件
func WriteFile_(filename string, data string) {
	if err := ioutil.WriteFile(filename, []byte(data), 0777); err != nil {
		panic(err)
	}
}
