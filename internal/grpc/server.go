package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"coordinateStoradge/internal/app"
	pb "coordinateStoradge/internal/grpc/api" // "github.com/Haba1234/coordinateStoradge/internal/grpc/api"
	"coordinateStoradge/internal/service"
	"google.golang.org/grpc"
)

//go:generate protoc -I ./api service.proto --go_out=. --go-grpc_out=.

// Server структура сервера.
type Server struct {
	mu *sync.Mutex
	pb.UnimplementedStatisticsServer
	srv    *grpc.Server
	search *service.Search
}

func NewServer(search *service.Search) *Server {
	return &Server{
		mu:     &sync.Mutex{},
		search: search,
	}
}

// Start запуск сервера gRPC.
func (s *Server) Start(addr string) error {
	log.Println("gRPC server " + addr + " running...")
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.srv = grpc.NewServer(grpc.StreamInterceptor(loggingServerInterceptor()))
	s.mu.Unlock()
	pb.RegisterStatisticsServer(s.srv, s)

	if err := s.srv.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func loggingServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println(fmt.Sprintf("method: %s, duration: %s, request: %+v", info.FullMethod, time.Since(time.Now()), srv))
		return handler(srv, ss)
	}
}

// Stop останов сервера gRPC.
func (s *Server) Stop() error {
	log.Println("gRPC server stopped")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.srv.GracefulStop()
	return nil
}

func (s *Server) StreamDots(stream pb.Statistics_StreamDotsServer) error {
	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("Клиент отключился")
			return nil
		}
		if err != nil {
			return err
		}

		point := app.Point{
			X: in.Point.GetX(),
			Y: in.Point.GetY(),
		}
		// Запрос соседей точки?
		if in.GetRequest() {
			if err := stream.Send(s.writeResponse(stream.Context(), point)); err != nil {
				return err
			}
			continue
		}

		// Запись точки в архив.
		s.search.SavePoint(point)
	}
}

// Подготовка данных к отправке.
func (s *Server) writeResponse(ctx context.Context, point app.Point) *pb.ServerStream {
	points := s.search.SearchNeighbors(ctx, point)
	pbPoints := make([]*pb.Result, app.MaxPoint)

	inPoint := pb.Point{
		X: point.X,
		Y: point.Y,
	}

	for i, p := range points {
		result := pb.Result{
			Point: &inPoint,
			Id:    int32(i + 1),
			Neighbor: &pb.Point{
				X: p.X,
				Y: p.Y,
			},
		}
		pbPoints[i] = &result
	}

	return &pb.ServerStream{
		Points: pbPoints,
	}
}
