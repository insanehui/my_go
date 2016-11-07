package gotcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

type MySvr struct{}

type MyPacket struct{ buff []byte }

// 实现Packet接口
func (this *MyPacket) Serialize() []byte {
	return this.buff
}

func NewPacket(buff []byte, hasLengthField bool) *MyPacket {
	p := &MyPacket{}

	if hasLengthField {
		p.buff = buff

	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff[0:4], uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

type Reader struct{}

func (this *Reader) ReadPacket(conn *net.TCPConn) (Packet, error) {
	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}
	if length = binary.BigEndian.Uint32(lengthBytes); length > 1024 {
		return nil, errors.New("the size of packet is larger than the limit")
	}

	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return nil, err
	}

	return NewPacket(buff, true), nil
}

func (this *MySvr) OnConnect(c *Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *MySvr) OnMessage(c *Conn, p Packet) bool {
	echoPacket := p.(*MyPacket)
	log.Printf("OnMessage:%+v", echoPacket)
	c.AsyncWritePacket(NewPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (this *MySvr) OnClose(c *Conn) {
	log.Println("OnClose:", c.GetExtraData())
}

func TestSimpleSvr(t *testing.T) {

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":1234")
	listener, _ := net.ListenTCP("tcp", tcpAddr)

	config := &Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := NewServer(config,
		&MySvr{},  // 回调接口
		&Reader{}) // 协议接口

	go srv.Start(listener, time.Second)
	time.Sleep(time.Second * 1)
	fmt.Println("listening:", listener.Addr())

	select {}

}
