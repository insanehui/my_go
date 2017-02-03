package gotcp

// 连接模块

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// 错误类型
var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

// 连接
type Conn struct {
	svr               *Server
	conn              *net.TCPConn  // the raw connection
	extraData         interface{}   // to save extra data
	closeOnce         sync.Once     // close the conn, once, per instance
	closeFlag         int32         // close flag
	closeChan         chan struct{} // close chanel
	packetSendChan    chan Packet   // packet send chanel
	packetReceiveChan chan Packet   // packeet receive chanel
}

// 工厂
func newConn(conn *net.TCPConn, srv *Server) *Conn {
	var ns, nr uint32
	if srv.config != nil {
		ns = srv.config.PacketSendChanLimit
		nr = srv.config.PacketReceiveChanLimit
	}
	if ns == 0 {
		ns = 20
	}
	if nr == 0 {
		nr = 20
	}

	return &Conn{
		svr:               srv,
		conn:              conn,
		closeChan:         make(chan struct{}),
		packetSendChan:    make(chan Packet, ns),
		packetReceiveChan: make(chan Packet, nr),
	}
}

func (c *Conn) GetExtraData() interface{} {
	return c.extraData
}

func (c *Conn) PutExtraData(data interface{}) {
	c.extraData = data
}

func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

// 关闭连接
func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		close(c.packetSendChan)
		close(c.packetReceiveChan)
		c.conn.Close()
		c.svr.callback.OnClose(c) // close 回调
	})
}

func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

// 异步写，不阻塞
func (c *Conn) AsyncWritePacket(p Packet, timeout time.Duration) (err error) {
	if c.IsClosed() {
		return ErrConnClosing
	}

	defer func() {
		if e := recover(); e != nil {
			err = ErrConnClosing
		}
	}()

	if timeout == 0 {
		select {
		case c.packetSendChan <- p:
			return nil

		default:
			return ErrWriteBlocking
		}

	} else {
		select {
		case c.packetSendChan <- p:
			return nil

		case <-c.closeChan:
			return ErrConnClosing

		case <-time.After(timeout):
			return ErrWriteBlocking
		}
	}
}

// 连接开始的主函数
func (c *Conn) Do() {

	if !c.svr.callback.OnConnect(c) { // connect 回调
		return
	}

	asyncDo(c.handleLoop, c.svr.waitGroup) // 所有异步协程的同步，都交给svr的wg来处理
	asyncDo(c.readLoop, c.svr.waitGroup)
	asyncDo(c.writeLoop, c.svr.waitGroup)
}

// 读流程
func (c *Conn) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.svr.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}

		p, err := c.svr.protocol.ReadPacket(c.conn)
		if err != nil {
			return
		}

		// 从网络读消息包，并存入读通道
		c.packetReceiveChan <- p
	}
}

// 写流程
func (c *Conn) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.svr.exitChan:
			return

		case <-c.closeChan:
			return

		// 从写通道取数据
		case p := <-c.packetSendChan:
			if c.IsClosed() {
				return
			}
			// 写入网络
			if _, err := c.conn.Write(p.Serialize()); err != nil {
				return
			}
		}
	}
}

// 从读通道取消息包，并调用 OnMessage 回调
func (c *Conn) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.svr.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.packetReceiveChan:
			if c.IsClosed() {
				return
			}
			if !c.svr.callback.OnMessage(c, p) {
				return
			}
		}
	}
}

func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
