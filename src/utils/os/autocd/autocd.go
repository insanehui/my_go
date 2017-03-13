// 自动cd到当前可执行文件的当前目录
package autocd

import (
	"os"
	"path/filepath"
)

func init() {
	dir := filepath.Dir(os.Args[0])
	os.Chdir(dir)
}

