package recover

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/logger"
	"runtime"
)

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
	logger.Log.Errorf("==> %s", string(buf[:n]))
}

func Recover() {
	if err := recover(); err != nil {
		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		PrintStack()
	}
}
