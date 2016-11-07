package io

import (
	"io/ioutil"
	"log"
)

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Panicf("read file err: %+v", err.Error())
	}

	return data
}
