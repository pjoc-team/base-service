package main

import (
	"fmt"
	"github.com/pjoc-team/base-service/pkg/util"
)

func main() {
	ip := util.GetHostIP()
	fmt.Println(ip)
}
