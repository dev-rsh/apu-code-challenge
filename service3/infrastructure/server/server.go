package server

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"service3/config"
	"services-challenge/proto/service3"
)

func Run() {
	lis, err := net.Listen("tcp", config.GetConfig().Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service3Server, err := newServer(config.GetConfig().Service1Host, config.GetConfig().Service2Host)
	if err != nil {
		log.Fatalf("failed to create service 3: %v", err)
	}
	s := grpc.NewServer()
	service3.RegisterService3Server(s, service3Server)
	log.Printf("Service 3 listening on %s", config.GetConfig().Host)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
