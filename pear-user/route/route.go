package router

import (
	"log"
	"net"

	"cyan.com/pear-common/discovery"
	"cyan.com/pear-common/logs"
	"cyan.com/pear-user/config"
	login_service_v1 "cyan.com/pear-user/pkg/service/login.service.v1"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

// Register Item/Objective
type Router interface {
	Register(r *gin.Engine) 
}


// 注册 center 方式1 在自身初始化
// 这种方式缺点是 得需要手动 _ 引入，方式2直接依赖
var Routers []Router
func Register(router ...Router)  {
	Routers = append(Routers, router...)
}

type GRPCConfig struct {
	Addr string
	RegisterFunc func(*grpc.Server)
}

func RegisterGRPC() *grpc.Server {
	c := GRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(s *grpc.Server) {
			login_service_v1.RegisterLoginServiceServer(s, login_service_v1.New())
		},
	}

	s := grpc.NewServer()
	c.RegisterFunc(s)

	lis, err := net.Listen("tcp", c.Addr)

	if err != nil {
		log.Println("cannot listen")
	}

	go func ()  {
		log.Printf("[INFO] GRPC Server listening on %s", c.Addr)
		err := s.Serve(lis)
		if err != nil {
			log.Println(err)
			return
		}
	}()
	return s
}

func InitRouter(r *gin.Engine) {
	for _, router := range Routers {
		router.Register(r)
	}
}


// 注册 center 方式 2, 在 center 初始化
// type RegisterRouter struct {}
// func New() *RegisterRouter {
// 	return &RegisterRouter{}
// }

// func (*RegisterRouter) Register(r *gin.Engine) {
// 	userRouter := &user.RouterUser{}
// 	userRouter.Register(r)
// }

func RegisterEtcdServer () *discovery.Register{
	etcdRegister := discovery.NewResolver(config.C.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	info := discovery.Server{
		Name: config.C.GC.Name,
		Addr: config.C.GC.Addr,
		Weight: config.C.GC.Weight,
		Version: config.C.GC.Version,
	}

	r := discovery.NewRegister(config.C.EC.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatal(err)
	}
	return r
}