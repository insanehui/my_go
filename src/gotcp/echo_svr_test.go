package gotcp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

type Callback struct{}

type EchoPacket struct{ buff []byte }

func (this *EchoPacket) Serialize() []byte {
	return this.buff
}

func (this *EchoPacket) GetLength() uint32 {
	return binary.BigEndian.Uint32(this.buff[0:4])
}

func (this *EchoPacket) GetBody() []byte {
	return this.buff[4:]
}

func NewEchoPacket(buff []byte, hasLengthField bool) *EchoPacket {
	p := &EchoPacket{}

	if hasLengthField {
		p.buff = buff

	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff[0:4], uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

type EchoProtocol struct{}

func (this *EchoProtocol) ReadPacket(conn *net.TCPConn) (Packet, error) {
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

	return NewEchoPacket(buff, true), nil
}

func (this *Callback) OnConnect(c *Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *Conn, p Packet) bool {
	echoPacket := p.(*EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (this *Callback) OnClose(c *Conn) {
	log.Println("OnClose:", c.GetExtraData())
}

func TestSvr(t *testing.T) {

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":8989")
	listener, _ := net.ListenTCP("tcp", tcpAddr)

	config := &Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := NewServer(config,
		&Callback{},     // 回调接口
		&EchoProtocol{}) // 协议接口

	go srv.Start(listener, time.Second)
	// go srv.Start(listener, time.Second)
	time.Sleep(time.Second * 2)
	fmt.Println("listening:", listener.Addr())

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)

	srv.Stop()
}

func TestClient(t *testing.T) {

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)

	echoProtocol := &EchoProtocol{}

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		conn.Write(NewEchoPacket([]byte("hello"), false).Serialize())

		// read
		p, err := echoProtocol.ReadPacket(conn)
		if err == nil {
			echoPacket := p.(*EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
		}

		time.Sleep(2 * time.Second)
	}

	conn.Close()

}
