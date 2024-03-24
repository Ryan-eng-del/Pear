package config

import (
	"log"
	"os"

	"cyan.com/pear-common/logs"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
	SC *ServerConfig
	GC *GrpcConfig
	EC *EtcdConfig
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
	conf.InitGrpcConfig()
	conf.InitEtcdConfig()

	return conf
}

type ServerConfig struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}


type EtcdConfig struct {
	Addrs []string `json:"addrs"`
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

func (c *Config) InitRedis() *redis.Options {
	option, _ := redis.ParseURL(c.viper.GetString("redis.source"))
	return option
}

var C = InitConfig()


type GrpcConfig struct {
	Name string
	Addr string
	Version string
	Weight int64
}

func (c *Config) InitGrpcConfig() {
	sc := &GrpcConfig{}
	sc.Name = c.viper.GetString("grpc.name")
	sc.Addr = c.viper.GetString("grpc.addr")
	sc.Version = c.viper.GetString("grpc.version")
	sc.Weight = c.viper.GetInt64("grpc.weight")
	c.GC = sc
}

func (c *Config) InitEtcdConfig() {
	sc := &EtcdConfig{}
	sc.Addrs = c.viper.GetStringSlice("etcd.addrs")
	c.EC = sc
}
