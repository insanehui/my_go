package gotcp

import (
	"net"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}

// 回调接口声明
type ConnCallback interface {

	// 返回false则关闭连接
	OnConnect(*Conn) bool

	// 返回false则关闭连接
	OnMessage(*Conn, Packet) bool

	OnClose(*Conn)
}

type Logic interface {
	Protocol
	ConnCallback
}
