package http

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func RProxy(url_path string) http.HandlerFunc {
	// TODO: 后续再支持一些灵活的批量代理策略（比如根据特定前缀），结合httprouter的特性应该可以做得到

	target, _ := url.Parse(url_path)

	return func(w http.ResponseWriter, r *http.Request) {

		o := new(http.Request)

		*o = *r

		// 重设host
		o.Host = target.Host

		// 重设url
		o.URL.Scheme = target.Scheme
		o.URL.Host = target.Host
		o.URL.Path = target.Path
		o.URL.RawQuery = r.URL.RawQuery

		// ======================= [COPIED BEGIN 以下代码拷贝而来未改动] ==================
		o.Proto = "HTTP/1.1"
		o.ProtoMajor = 1
		o.ProtoMinor = 1
		o.Close = false

		transport := http.DefaultTransport

		res, err := transport.RoundTrip(o)

		if err != nil {
			log.Printf("http: proxy error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hdr := w.Header()

		for k, vv := range res.Header {
			for _, v := range vv {
				hdr.Add(k, v)
			}
		}

		for _, c := range res.Cookies() {
			w.Header().Add("Set-Cookie", c.Raw)
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			io.Copy(w, res.Body)
		}
		// ======================= [COPIED END] ==================
	}
}

// ====================== 以下是网上摘抄的原始代码，供今后比较参考 ===========================
var targetURL = &url.URL{
	Host: "www.baidu.com",
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func handler(w http.ResponseWriter, r *http.Request) {
	o := new(http.Request)

	*o = *r

	o.Host = targetURL.Host
	o.URL.Scheme = targetURL.Scheme
	o.URL.Host = targetURL.Host
	o.URL.Path = singleJoiningSlash(targetURL.Path, o.URL.Path)

	if q := o.URL.RawQuery; q != "" {
		o.URL.RawPath = o.URL.Path + "?" + q
	} else {
		o.URL.RawPath = o.URL.Path
	}

	o.URL.RawQuery = targetURL.RawQuery

	o.Proto = "HTTP/1.1"
	o.ProtoMajor = 1
	o.ProtoMinor = 1
	o.Close = false

	transport := http.DefaultTransport

	res, err := transport.RoundTrip(o)

	if err != nil {
		log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hdr := w.Header()

	for k, vv := range res.Header {
		for _, v := range vv {
			hdr.Add(k, v)
		}
	}

	for _, c := range res.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(w, res.Body)
	}
}

// func main() {
// 	http.HandleFunc("/", handler)
// 	http.ListenAndServe(":1234", nil)
// }
