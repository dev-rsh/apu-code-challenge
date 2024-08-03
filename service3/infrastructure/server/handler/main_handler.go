package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"service3/dto"
	"service3/infrastructure/repository"
	"services-challenge/proto/service1"
	"services-challenge/proto/service2"
	"services-challenge/proto/service3"
	"time"
)

type Server struct {
	service3.UnimplementedService3Server
	Service1Client service1.Service1Client
	Service2Client service2.Service2Client
}

func (s *Server) GetResult(ctx context.Context, in *service3.Empty) (*service3.ResultResponse, error) {
	svc1Res := make(chan *service1.DataResponse)
	svc1Err := make(chan error)
	svc1Delay := make(chan time.Duration)
	svc2Res := make(chan *service2.DataResponse)
	svc2Err := make(chan error)
	svc2Delay := make(chan time.Duration)
	go s.getResultFromService1(svc1Res, svc1Err, svc1Delay)
	go s.getResultFromService2(svc2Res, svc2Err, svc2Delay)

	var (
		service1Res, service2Res         string
		service1Success, service2Success bool
		service1Delay, service2Delay     time.Duration
	)

	select {
	case data := <-svc1Res:
		service1Res = data.Data
		service1Success = true
	case err := <-svc1Err:
		if status.Code(err) == codes.DeadlineExceeded {
			service1Res = "Timeout exceeded"
		} else {
			service1Res = err.Error()
		}
	}
	service1Delay = <-svc1Delay

	select {
	case data := <-svc2Res:
		service2Res = data.Data
		service2Success = true
	case err := <-svc2Err:
		if status.Code(err) == codes.DeadlineExceeded {
			service2Res = "Timeout exceeded"
		} else {
			service2Res = err.Error()
		}
	}
	service2Delay = <-svc2Delay

	go func() {
		serviceResultRepository := repository.GetServiceResultRepository()
		saveResultToDb(serviceResultRepository, service1Success, service2Success, service1Delay, service2Delay)
	}()

	return &service3.ResultResponse{
		Service1Data: service1Res,
		Service2Data: service2Res,
	}, nil
}

func (s *Server) getResultFromService1(resultChan chan<- *service1.DataResponse, errChan chan<- error, delayChan chan<- time.Duration) {
	timeNow := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	data, err := s.Service1Client.GetData(ctx, &service1.Empty{})
	delay := time.Since(timeNow)
	if err != nil {
		errChan <- err
		delayChan <- delay
		return
	}

	resultChan <- data
	delayChan <- delay
}

func (s *Server) getResultFromService2(resultChan chan<- *service2.DataResponse, errChan chan<- error, delayChan chan<- time.Duration) {
	timeNow := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	data, err := s.Service2Client.GetData(ctx, &service2.Empty{})
	delay := time.Since(timeNow)
	if err != nil {
		errChan <- err
		delayChan <- delay
		return
	}

	resultChan <- data
	delayChan <- delay
}

func saveResultToDb(repo repository.ServiceResultRepository, svc1Success, svc2Success bool, svc1Delay, svc2Delay time.Duration) {
	svcResDto := dto.ServiceResultDto{
		ExecutionTime:   time.Now(),
		Service1Success: svc1Success,
		Service2Success: svc2Success,
		Service1Delay:   svc1Delay.Milliseconds(),
		Service2Delay:   svc2Delay.Milliseconds(),
	}

	err := repo.Save(svcResDto)
	if err != nil {
		log.Printf("There was a problem writing results to db. Err is %s", err)
	}
}
