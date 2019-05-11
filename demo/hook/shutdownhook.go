package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errc := make(chan error)

	go func() {
		log.Println("Listening signals...")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("Signal %v", <-c)
	}()

	log.Println("Exit:", <-errc)
}
