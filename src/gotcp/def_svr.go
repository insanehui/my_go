package gotcp

// 最基础（缺省）的服务器逻辑, 预置了缺省的回调

import (
	"log"
)

type Def struct{}

func (this *Def) OnConnect(c *Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	log.Println("OnConnect:", addr)
	return true
}

func (this *Def) OnMessage(c *Conn, p Packet) bool {
	addr := c.GetRawConn().RemoteAddr()
	log.Printf("OnMessage:%+v:%+v", addr, p)
	return true
}

func (this *Def) OnClose(c *Conn) {
	addr := c.GetRawConn().RemoteAddr()
	log.Printf("OnClose:%+v", addr)
}
