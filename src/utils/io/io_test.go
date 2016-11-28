package io

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func TestReadFile(t *testing.T) {
	data := ReadFile_("io_test.go")
	log.Printf("%+v", string(data))
}

func Test_logwriter(t *testing.T) {
	b := new(bytes.Buffer)
	l := LogWriter(b)
	io.WriteString(l, `this is a test string
hahaha
shit safd 
asfasfs
`)
	log.Printf("%s", b)
}
