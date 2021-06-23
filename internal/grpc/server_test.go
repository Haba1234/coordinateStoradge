package grpc

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/Haba1234/coordinateStoradge/internal/app"
	pb "github.com/Haba1234/coordinateStoradge/internal/grpc/api"
	"github.com/Haba1234/coordinateStoradge/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Suite struct {
	suite.Suite
	service    *mocks.Search
	stream     *mocks.Statistics_StreamDotsServer
	ind        int
	testPoints []pb.Point
	points     []app.Point
	errorFn    func() error
	streamFn   func() *pb.ClientStream
}

func (s *Suite) SetupSuite() {
	s.testPoints = []pb.Point{
		{
			X: 1000,
			Y: 1000,
		},
		{
			X: 2000,
			Y: 2000,
		},
		{
			X: 3000,
			Y: 3000,
		},
	}

	s.points = []app.Point{
		{
			X: 1000,
			Y: 1000,
		},
		{
			X: 2000,
			Y: 2000,
		},
		{
			X: 3000,
			Y: 3000,
		},
	}

	s.errorFn = func() error {
		if s.ind > len(s.testPoints) {
			return io.EOF
		}
		return nil
	}

	s.streamFn = func() *pb.ClientStream {
		i := s.ind
		s.ind = i + 1
		flag := false
		if i == len(s.testPoints)-1 {
			flag = true
		}
		if s.ind > len(s.testPoints) {
			return nil
		}
		return &pb.ClientStream{
			Point:   &s.testPoints[i],
			Request: flag,
		}
	}
}

func (s *Suite) SetupTest() {
	s.service = new(mocks.Search)
	s.stream = new(mocks.Statistics_StreamDotsServer)
}

func (s *Suite) TearDownTest() {
	s.service.AssertExpectations(s.T())
	s.stream.AssertExpectations(s.T())
}

func (s *Suite) TestStartStop() {
	defer goleak.VerifyNone(s.T())
	server := NewServer(s.service)
	go func() {
		err := server.Start(":8080")
		s.NoError(err)
	}()

	time.Sleep(10 * time.Millisecond)
	err := server.Stop()
	s.NoError(err)
}

func (s *Suite) TestStreamDotsOK() {
	s.stream.On("Recv").Return(s.streamFn, s.errorFn)
	s.stream.On("Context").Return(context.Background())
	s.stream.On("Send", mock.Anything).Return(nil)
	s.service.On("SavePoint", mock.Anything).Return()
	s.service.On("SearchNeighbors", context.Background(), mock.Anything).Return(s.points)

	server := NewServer(s.service)
	err := server.StreamDots(s.stream)
	s.NoError(err)
}

func (s *Suite) TestStreamDotsError() {
	s.stream.On("Recv").Return(nil, status.Error(codes.Aborted, "client aborted"))
	server := NewServer(s.service)
	err := server.StreamDots(s.stream)
	s.Error(err)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
