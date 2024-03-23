package config

import (
	"log"
	"os"

	"cyan.com/pear-common/logs"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
	SC *ServerConfig
}

func InitConfig () *Config {
	conf := &Config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("app")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	log.Println(workDir)
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}


	conf.InitServer()
	conf.InitZapLog()

	return conf
}

type ServerConfig struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

func (c *Config) InitServer() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}


func (c *Config) InitZapLog() {
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName: c.viper.GetString("zap.infoFileName"),
		WarningFileName: c.viper.GetString("zap.warnFileName"),
		MaxSize: c.viper.GetInt("zap.maxSize"), // 10M
		MaxAge: c.viper.GetInt("zap.maxAge"), // 24hour
		MaxBackups: c.viper.GetInt("zap.maxBackups"), // 6 backups
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatal(err)
	}
}


var C = InitConfig()

