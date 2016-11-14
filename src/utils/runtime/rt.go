package runtime

import (
	"log"
	"runtime"
)

// 取得运行时栈（调用堆栈）
func Stack() string {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return string(buf)
}

// 将调用堆栈打日志
func Log() {
	log.Printf("################## ERR!! STACK BEGIN .... ###################")
	log.Printf(Stack())
	log.Printf("########################## STACK END ########################")
}
