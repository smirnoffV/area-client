package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	pb "github.com/smirnoffV/area-client/pb"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"time"
)

var (
	serverAddress = flag.String("addr", ":10000", "")
	minNumber     = flag.Int("min", 0, "")
	maxNumber     = flag.Int("max", 50000, "")
)

func main() {
	log.Info("starting gRPC client")

	flag.Parse()

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	areaClient := pb.NewAreaClient(conn)

	res, err := areaClient.Circle(context.Background(), &pb.CircleRequest{Radius: 3})
	if err != nil {
		panic(err)
	}

	log.Printf("Circle area = %v", res.Area)

	res, err = areaClient.Square(context.Background(), &pb.SquareRequest{Side: 5})
	if err != nil {
		panic(err)
	}

	log.Printf("Square area = %v", res.Area)

	res, err = areaClient.Rectangle(context.Background(), &pb.RectangleRequest{Width: 10, Height: 3})
	if err != nil {
		panic(err)
	}

	log.Printf("Rectangle area = %v", res.Area)

	ctx, cancelFn := context.WithCancel(context.Background())
	streamMaxClient, err := areaClient.Max(ctx)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-time.NewTicker(time.Millisecond * 500).C:
				newNumber := int64(random(*minNumber, *maxNumber))
				log.Infof("Send new number: %v", newNumber)
				if err := streamMaxClient.Send(&pb.NumberRequest{Number: newNumber}); err != nil {
					log.WithError(err).Error("error publishing message into stream")
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer cancelFn()

		for range time.NewTimer(time.Minute * 5).C {
			log.Info("stop publishing messages to the stream")
			return
		}
	}()

	wg.Wait()
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
