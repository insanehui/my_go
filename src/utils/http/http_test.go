package http

import (
	"testing"
	"log"
)

func Test_map2urlvals(t *testing.T) {
	m := map[string]interface{}{
		"aa" : 1,
		"bb" : "hahah",
		"c" : 334.5 }
	r := toUrlVals(m)
	log.Printf("r: %+v", r)
}

func Test_post(t *testing.T) {
	res, _ := Post("http://www.baidu.com", nil)
	log.Println(res)
}
