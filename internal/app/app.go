package app

import (
	"context"

	pb "github.com/Haba1234/coordinateStoradge/internal/grpc/api"
)

type Point struct {
	X, Y uint32
}

const (
	MaxPoint = 3     // Кол-во искомых соседей
	MaxLimit = 10000 // Максимально допустимая координатна точки
)

type Storage interface {
	ReadPoint(z uint64) (bool, bool)
	RLock()
	RUnlock()
	Len() int
	AddPoint(z uint64)
}

type Search interface {
	SavePoint(point Point)
	SearchNeighbors(ctx context.Context, point Point) []Point
}

type Server interface {
	Start(addr string) error
	Stop() error
	StreamDots(stream pb.Statistics_StreamDotsServer) error
}
