package handler

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"service3/dto"
	"services-challenge/proto/service1"
	"services-challenge/proto/service2"
	"services-challenge/proto/service3"
	"testing"
	"time"
)

type MockService1Client struct {
	mock.Mock
}

func (m *MockService1Client) GetData(ctx context.Context, in *service1.Empty, opts ...grpc.CallOption) (*service1.DataResponse, error) {
	args := m.Called(ctx, in)
	var resp *service1.DataResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*service1.DataResponse)
	} else {
		resp = nil
	}
	return resp, args.Error(1)
}

type MockService2Client struct {
	mock.Mock
}

func (m *MockService2Client) GetData(ctx context.Context, in *service2.Empty, opts ...grpc.CallOption) (*service2.DataResponse, error) {
	args := m.Called(ctx, in)
	var resp *service2.DataResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*service2.DataResponse)
	} else {
		resp = nil
	}
	return resp, args.Error(1)
}

func TestGetResult(t *testing.T) {

	t.Run("Successful Response from Both Services", func(t *testing.T) {
		mockService1Client := new(MockService1Client)
		mockService2Client := new(MockService2Client)
		server := &Server{
			Service1Client: mockService1Client,
			Service2Client: mockService2Client,
		}

		mockService1Client.On("GetData", mock.Anything, &service1.Empty{}, mock.Anything).Return(&service1.DataResponse{Data: "Data from Service 1"}, nil)
		mockService2Client.On("GetData", mock.Anything, &service2.Empty{}, mock.Anything).Return(&service2.DataResponse{Data: "Data from Service 2"}, nil)

		resp, err := server.GetResult(context.Background(), &service3.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "Data from Service 1", resp.Service1Data)
		assert.Equal(t, "Data from Service 2", resp.Service2Data)
	})

	t.Run("Timeout from Service 1", func(t *testing.T) {
		mockService1Client := new(MockService1Client)
		mockService2Client := new(MockService2Client)
		server := &Server{
			Service1Client: mockService1Client,
			Service2Client: mockService2Client,
		}

		mockService1Client.On("GetData", mock.Anything, &service1.Empty{}, mock.Anything).Return(nil, status.Error(codes.DeadlineExceeded, "timeout"))
		mockService2Client.On("GetData", mock.Anything, &service2.Empty{}, mock.Anything).Return(&service2.DataResponse{Data: "Data from Service 2"}, nil)

		resp, err := server.GetResult(context.Background(), &service3.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "Timeout exceeded", resp.Service1Data)
		assert.Equal(t, "Data from Service 2", resp.Service2Data)
	})

	t.Run("Timeout from Service 2", func(t *testing.T) {
		mockService1Client := new(MockService1Client)
		mockService2Client := new(MockService2Client)
		server := &Server{
			Service1Client: mockService1Client,
			Service2Client: mockService2Client,
		}

		mockService1Client.On("GetData", mock.Anything, &service1.Empty{}, mock.Anything).Return(&service1.DataResponse{Data: "Data from Service 1"}, nil)
		mockService2Client.On("GetData", mock.Anything, &service2.Empty{}, mock.Anything).Return(nil, status.Error(codes.DeadlineExceeded, "timeout"))

		resp, err := server.GetResult(context.Background(), &service3.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "Data from Service 1", resp.Service1Data)
		assert.Equal(t, "Timeout exceeded", resp.Service2Data)
	})

	t.Run("Error from Service 1", func(t *testing.T) {
		mockService1Client := new(MockService1Client)
		mockService2Client := new(MockService2Client)
		server := &Server{
			Service1Client: mockService1Client,
			Service2Client: mockService2Client,
		}

		mockService1Client.On("GetData", mock.Anything, &service1.Empty{}, mock.Anything).Return(nil, errors.New("Service 1 error"))
		mockService2Client.On("GetData", mock.Anything, &service2.Empty{}, mock.Anything).Return(&service2.DataResponse{Data: "Data from Service 2"}, nil)

		resp, err := server.GetResult(context.Background(), &service3.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "Service 1 error", resp.Service1Data)
		assert.Equal(t, "Data from Service 2", resp.Service2Data)
	})

	t.Run("Error from Service 2", func(t *testing.T) {
		mockService1Client := new(MockService1Client)
		mockService2Client := new(MockService2Client)
		server := &Server{
			Service1Client: mockService1Client,
			Service2Client: mockService2Client,
		}

		mockService1Client.On("GetData", mock.Anything, &service1.Empty{}, mock.Anything).Return(&service1.DataResponse{Data: "Data from Service 1"}, nil)
		mockService2Client.On("GetData", mock.Anything, &service2.Empty{}, mock.Anything).Return(nil, errors.New("Service 2 error"))

		resp, err := server.GetResult(context.Background(), &service3.Empty{})
		assert.NoError(t, err)
		assert.Equal(t, "Data from Service 1", resp.Service1Data)
		assert.Equal(t, "Service 2 error", resp.Service2Data)
	})

}

type MockServiceResultRepository struct {
	mock.Mock
}

func (m *MockServiceResultRepository) Save(dto dto.ServiceResultDto) error {
	args := m.Called(dto)
	return args.Error(0)
}

func TestSaveResultToDb(t *testing.T) {
	mockRepo := new(MockServiceResultRepository)

	testCases := []struct {
		name          string
		svc1Success   bool
		svc2Success   bool
		svc1Delay     time.Duration
		svc2Delay     time.Duration
		expectedError bool
		repositoryErr error
	}{
		{
			name:          "Successful save",
			svc1Success:   true,
			svc2Success:   true,
			svc1Delay:     100 * time.Millisecond,
			svc2Delay:     200 * time.Millisecond,
			expectedError: false,
			repositoryErr: nil,
		},
		{
			name:          "Repository error",
			svc1Success:   false,
			svc2Success:   true,
			svc1Delay:     300 * time.Millisecond,
			svc2Delay:     150 * time.Millisecond,
			expectedError: true,
			repositoryErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("Save", mock.MatchedBy(func(dto dto.ServiceResultDto) bool {
				return dto.Service1Success == tc.svc1Success &&
					dto.Service2Success == tc.svc2Success &&
					dto.Service1Delay == tc.svc1Delay.Milliseconds() &&
					dto.Service2Delay == tc.svc2Delay.Milliseconds()
			})).Return(tc.repositoryErr)

			saveResultToDb(mockRepo, tc.svc1Success, tc.svc2Success, tc.svc1Delay, tc.svc2Delay)

			mockRepo.AssertExpectations(t)

			if tc.expectedError {
				mockRepo.AssertCalled(t, "Save", mock.Anything)
			}
		})
	}
}
