package main

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand"
	"time"

	pb "coordinateStoradge/internal/grpc/api"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Client started")
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client start error: %v", err)
	}
	defer conn.Close()

	client := pb.NewStatisticsClient(conn)
	stream, err := client.StreamDots(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Println("Принято:", in.String())
		}
	}()

	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	pointCount := 1000000
	points := make([]*pb.Point, pointCount)

	for i := 0; i < pointCount; i++ {
		points[i] = randomPoint(r)
	}

	var oldPoint *pb.Point
	for i, point := range points {
		request := false
		if i > 0 && i%10000 == 0 { // Каждый 10 000 -й вызов, это запрос соседей.
			request = true
			point = oldPoint
		}

		if err := stream.Send(&pb.ClientStream{
			Point:   point,
			Request: request,
		}); err != nil {
			log.Panicf("Failed to send a note: %v", err)
		}
		oldPoint = point
	}
	stream.CloseSend()
	<-waitc
	log.Println("Client stopped")
}

func randomPoint(r *rand.Rand) *pb.Point {
	return &pb.Point{
		X: uint32(r.Intn(9999) + 1),
		Y: uint32(r.Intn(9999) + 1),
	}
}
