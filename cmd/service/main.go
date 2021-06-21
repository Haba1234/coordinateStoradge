package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os/signal"
	"syscall"

	"coordinateStoradge/internal/grpc"
	"coordinateStoradge/internal/service"
	"coordinateStoradge/internal/storage"
)

func main() {
	var portServer string
	flag.StringVar(&portServer, "port", "8080", "gRPC server port number")
	flag.Parse()

	newStorage := storage.NewStorage()
	search := service.NewSearch(newStorage)
	server := grpc.NewServer(search)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	addrServer := net.JoinHostPort("", portServer)
	go func() {
		<-ctx.Done()

		if err := server.Stop(); err != nil {
			log.Println("failed to stop gRPC server: " + err.Error())
		}
	}()

	log.Println("Start service")
	if err := server.Start(addrServer); err != nil {
		log.Println("failed to start gRPC server: " + err.Error())
		cancel()
	}
}
