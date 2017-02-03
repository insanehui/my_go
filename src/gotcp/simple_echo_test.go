package gotcp

import (
	// "errors"
	// "io"
	"log"
	"net"
	"testing"
	"time"

	I "utils/io"
)

// =============== 结构体声明 =======================

type MyPacket struct{ buff string }

// 必须实现Packet接口
func (this *MyPacket) Serialize() []byte {
	return []byte(this.buff)
}

type TheSvr struct {
	Def
}

func (this *TheSvr) ReadPacket(conn *net.TCPConn) (Packet, error) {

	b := make([]byte, 100)
	p := MyPacket{}
	// n, err := io.ReadAtLeast(conn, b, 10)
	n, err := I.ReadBy(conn, b, "\n")

	if err != nil {
		return nil, err
	}
	log.Printf("read %+v bytes: %+v", n, string(b))
	p.buff = string(b)
	return &p, nil
}

func (this *TheSvr) OnMessage(c *Conn, p Packet) bool {
	log.Printf("packet: %+v", p)
	c.AsyncWritePacket(p, time.Second)
	return true
}

func TestSimpleSvr(t *testing.T) {

	svr := New(&TheSvr{}, &Config{Addr: ":2222"})
	go svr.Run()
	log.Println("svr running")

	select {}

}
