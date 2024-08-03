package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"service1/config"
	"services-challenge/proto/service1"
	"time"
)

type server struct {
	service1.UnimplementedService1Server
}

func (s *server) GetData(ctx context.Context, in *service1.Empty) (*service1.DataResponse, error) {
	randGen := rand.New(rand.NewSource(time.Now().UnixMilli()))
	randSleep := randGen.Intn(600-300) + 300
	time.Sleep(time.Duration(randSleep) * time.Millisecond)

	return &service1.DataResponse{Data: "Data from Service 1"}, nil
}

func main() {
	lis, err := net.Listen("tcp", config.GetConfig().Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service1.RegisterService1Server(s, &server{})
	log.Printf("Service 1 listening on %s", config.GetConfig().Host)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
