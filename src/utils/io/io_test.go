package io

import (
	"log"
	"testing"
)

func TestReadFile(t *testing.T) {
	data := ReadFile("io_test.go")
	log.Printf("%+v", string(data))
}
