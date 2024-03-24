package main

import (
	"os"
	"os/signal"
	"syscall"

	router "cyan.com/pear-user/route"
)


func main() {
	s := router.RegisterGRPC()
	quit := make(chan os.Signal, 1)
	r := router.RegisterEtcdServer()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit
	s.Stop()
	r.Stop()
	// common.Run(r, s, config.C.SC.Name, config.C.SC.Addr)
}