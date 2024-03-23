package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "cyan.com/pear-user/api"
	router "cyan.com/pear-user/route"
)


func main() {
	// r := gin.Default()
	// router.InitRouter(r)
	s := router.RegisterGRPC()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit
	s.Stop()
	// common.Run(r, s, config.C.SC.Name, config.C.SC.Addr)
}