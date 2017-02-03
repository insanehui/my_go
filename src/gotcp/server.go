package gotcp

// Server类

import (
	"log"
	"net"
	"sync"
)

type Config struct {
	PacketSendChanLimit    uint32 // 写缓冲长度（包的数目）
	PacketReceiveChanLimit uint32 // 读缓冲

	Addr string // 监听的地址, 如 ":1234"
}

type Server struct {
	config    *Config         // server configuration
	callback  ConnCallback    // message callbacks in connection
	protocol  Protocol        // customize packet protocol
	exitChan  chan struct{}   // 用于传递退出信号。（会停止所有的协程？）
	waitGroup *sync.WaitGroup // 等待所有的协程
}

// 工厂
func NewServer(config *Config, callback ConnCallback, protocol Protocol) *Server {
	return &Server{
		config:    config,
		callback:  callback,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

func New(logic Logic, config *Config) *Server {
	return &Server{
		config:    config,
		callback:  logic,
		protocol:  logic,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

// 启动服务
// 注：死循环。需要go出去使用
func (s *Server) Start(listener *net.TCPListener) {
	s.waitGroup.Add(1)

	defer func() {
		if p := recover(); p != nil {
			log.Printf("ERROR: %+v", p)
		}

		listener.Close()
		s.waitGroup.Done()
	}()

	// 主循环
	for {

		select {
		case <-s.exitChan: // 如果收到退出信号，则退出
			return

		default: // 否则继续服务
		}

		// accept
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("continue to listen...")
			continue
		}

		s.waitGroup.Add(1)
		go func() {
			newConn(conn, s).Do() // 处理新连接的入口主函数
			s.waitGroup.Done()
		}()
	}
}

// 启动服务的新版。在配置里设置监听端口
func (s *Server) Run() {
	s.waitGroup.Add(1)

	var addr string
	if s.config != nil {
		addr = s.config.Addr
	}
	if addr == "" {
		addr = ":1234"
	}

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)
	listener, _ := net.ListenTCP("tcp", tcpAddr)

	defer func() {
		if p := recover(); p != nil {
			log.Printf("ERROR: %+v", p)
		}

		listener.Close()
		s.waitGroup.Done()
	}()

	// 主循环
	for {

		select {
		case <-s.exitChan: // 如果收到退出信号，则退出
			return

		default: // 否则继续服务
		}

		// accept
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("continue to listen...")
			continue
		}

		s.waitGroup.Add(1)
		go func() {
			newConn(conn, s).Do() // 处理新连接的入口主函数
			s.waitGroup.Done()
		}()
	}
}

// 停止服务
func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait() // 等待所有的协程执行完毕
}
