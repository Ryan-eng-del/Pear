package user

import (
	"log"

	"cyan.com/pear-api/config"
	login_service_v1 "cyan.com/pear-api/pkg/service/login.service.v1"
	"cyan.com/pear-common/discovery"
	"cyan.com/pear-common/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var LoginServiceClient login_service_v1.LoginServiceClient

func InitUserClient() {
	etcdRegister := discovery.NewResolver(config.C.EC.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.Dial("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err  != nil {
		log.Fatalf("did not connect to server %v", err)
	}
	LoginServiceClient = login_service_v1.NewLoginServiceClient(conn)
}