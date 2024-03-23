package main

import (
	_ "cyan.com/pear-api/api"
	"cyan.com/pear-api/config"
	router "cyan.com/pear-api/route"
	common "cyan.com/pear-common"
	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()
	router.InitRouter(r)
	common.Run(r, nil, config.C.SC.Name, config.C.SC.Addr)
}