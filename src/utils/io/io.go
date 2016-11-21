package io

import (
	"io/ioutil"
	"log"
	"os"
)

// 判断文件是否存在
func Exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err) // 不明白这个IsExist是干什么的？
}

// 判断文件是否存在
func Exists_(p string) {
	_, err := os.Stat(p)
	if  err == nil || os.IsExist(err) {
		panic( p + " not exists")
	}
}

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
