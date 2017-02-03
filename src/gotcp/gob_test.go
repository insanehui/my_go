package gotcp

// 测试gob svr

import (
	"encoding/gob"
	"log"
	"net"
	"testing"
	"time"
)

type gobSvr struct {
	Gob
}

func TestGobSvr(t *testing.T) {

	svr := New(&gobSvr{}, nil)
	go svr.Run()
	log.Println("gob svr running")

	select {}

}

func TestGobClient(t *testing.T) {

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:1234")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)

	enc := gob.NewEncoder(conn)
	a := GobData{"a": 1, "b": GobData{"c": "shit"}}

	e := enc.Encode(a)
	if e != nil {
		log.Printf("error: %+v", e)
	}

	time.Sleep(2 * time.Second)
	conn.Close()
}
