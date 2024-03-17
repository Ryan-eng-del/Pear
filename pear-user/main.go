package main

import (
	common "cyan.com/pear-common"
	"cyan.com/pear-common/logs"
	_ "cyan.com/pear-user/api"
	router "cyan.com/pear-user/route"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func main() {
	r := gin.Default()
	router.InitRouter(r)
	lc := &logs.LogConfig{
		DebugFileName: "/Users/max/Documents/coding/Backend/Golang/Personal/Go-Pear-Project/pear-common/logs/debug.log",
		InfoFileName: "/Users/max/Documents/coding/Backend/Golang/Personal/Go-Pear-Project/pear-common/logs/info.log",
		WarningFileName: "/Users/max/Documents/coding/Backend/Golang/Personal/Go-Pear-Project/pear-common/logs/error.log",
		MaxSize: 10, // 10M
		MaxAge: 24, // 24hour
		MaxBackups: 3, // 6 backups
	}
	logs.InitLogger(lc)
	zap.L().Debug("Initializing")
	zap.L().Info("Initializing")
	zap.L().Warn("Initializing")
	zap.L().Error("Initializing")
	common.Run(r, "pear-project", ":8080")
}