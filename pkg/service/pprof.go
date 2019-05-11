package service

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	StartPProf()
}

func StartPProf() {
	StartPProfWithPort(6060)
}

func StartPProfWithPort(port int) {
	StartPProfWithAddr(fmt.Sprintf("localhost:%d", port))
}
func StartPProfWithAddr(addr string) {
	go func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		log.Println(http.ListenAndServe(addr, nil))
	}()
}
