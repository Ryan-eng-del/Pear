package common

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Run(r *gin.Engine, s *grpc.Server, srvName string, addr string) {
	srv := &http.Server{
		Addr: ":80",
		Handler: r,
	}
	go func(){
		log.Printf("web server listening on %s", srv.Addr)
		log.Fatalln(srv.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit
	log.Println("[INFO] Gracefully start shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] HTTPServerStop failed: %v\n", err)
	}

	if s != nil {
		s.Stop()
	}

	log.Println("[INFO] Gracefully end shutting down")
}