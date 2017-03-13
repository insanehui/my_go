package stdout

import (
	"log"
	"os"
)

// 将日志打到stdout，似乎go缺省是输出到了stderr
func init() {
	log.SetOutput(os.Stdout)
}
