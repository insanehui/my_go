package stdout

import (
	"log"
	"os"
)

// 将日志打到stdout
func init() {
	log.SetOutput(os.Stdout)
}
