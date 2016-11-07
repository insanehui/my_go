package utils

import (
	"log"
	"reflect"
	"testing"
)

func TestMd5(t *testing.T) {
	if Md5("1234") != "81dc9bdb52d04dc20036dbd8313ed055" {
		t.Errorf(`"1234"的md5不正确`)
	}
}

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "v": 1}
	keys := Keys(m)
	log.Printf("%+v", keys)
	if !reflect.DeepEqual(keys, []interface{}{"a", "b", "v"}) {
		t.Error()
	}
}

func dummy() {
	log.Printf("dummy")
}
