package gotcp

// 基于gob来通信的服务

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type Gob struct {
	Def
}

type GobData map[string]interface{}

type GobPkt struct {
	Data GobData
}

func init() {
	gob.Register(GobData{})
}

// 实现Packet接口
func (me *GobPkt) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	e := enc.Encode(me.Data)
	if e != nil {
		log.Printf("error: %+v", e)
	}
	log.Printf("encoded, len: %+v", buf.Len())
	return buf.Bytes()
}

func (this *Gob) ReadPacket(conn *net.TCPConn) (Packet, error) {

	p := GobPkt{}
	dec := gob.NewDecoder(conn)

	// 先暂时使用每次new 编解码器的方法，后续put到 conn 的extra data里去
	e := dec.Decode(&p.Data)
	if e != nil {
		return nil, e
	}
	log.Printf("decoded data: %+v", p)
	return &p, nil
}
