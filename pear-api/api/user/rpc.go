package user

import (
	"log"

	login_service_v1 "cyan.com/pear-api/pkg/service/login.service.v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var LoginServiceClient login_service_v1.LoginServiceClient

func InitUserClient() {
	conn, err := grpc.Dial("127.0.0.1:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err  != nil {
		log.Fatalf("did not connect to server %v", err)
	}
	LoginServiceClient = login_service_v1.NewLoginServiceClient(conn)
}