package server

import (
	"google.golang.org/grpc"
	"service3/infrastructure/server/handler"
	"services-challenge/proto/service1"
	"services-challenge/proto/service2"
)

func newServer(service1Addr, service2Addr string) (*handler.Server, error) {
	conn1, err := grpc.Dial(service1Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	conn2, err := grpc.Dial(service2Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &handler.Server{
		Service1Client: service1.NewService1Client(conn1),
		Service2Client: service2.NewService2Client(conn2),
	}, nil
}
