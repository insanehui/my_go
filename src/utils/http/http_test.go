package http

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	_ "utils/log/stdout"
)

func Test_map2urlvals(t *testing.T) {
	m := map[string]interface{}{
		"aa": 1,
		"bb": "hahah",
		"c":  334.5}
	r := toUrlVals(m)
	log.Printf("r: %+v", r)
}

func Test_post(t *testing.T) {
	res, _ := Post("http://www.baidu.com", nil)
	log.Println(res)
}

func Test_jsondo(t *testing.T) {
	type Para struct {
		Name string `valid:"-"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p Para

		var ret struct {
			Err
			Name string
		}
		JsonDo(w, r, &p, &ret, func() {
			ret.Name = p.Name
		})
	}))
	defer ts.Close()

	res, _ := Post(ts.URL, Para{"hello world"})
	log.Println("client get: ", res)

}
