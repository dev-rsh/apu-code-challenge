package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"service2/config"
	"services-challenge/proto/service2"
	"time"
)

type server struct {
	service2.UnimplementedService2Server
}

func (s *server) GetData(ctx context.Context, in *service2.Empty) (*service2.DataResponse, error) {
	randGen := rand.New(rand.NewSource(time.Now().UnixMilli()))
	randSleep := randGen.Intn(600-300) + 300
	time.Sleep(time.Duration(randSleep) * time.Millisecond)
	return &service2.DataResponse{Data: "Data from Service 2"}, nil
}

func main() {
	lis, err := net.Listen("tcp", config.GetConfig().Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service2.RegisterService2Server(s, &server{})
	log.Printf("Service 2 listening on %s", config.GetConfig().Host)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
