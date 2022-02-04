// Sample usage of the gocron-server pkg
package main

import (
	"fmt"
	"log"
	"net"

	gocron_server "github.com/strnophix/gocron-server/pkg"
	pb "github.com/strnophix/gocron-server/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9092"
)

type Counter struct {
	Current int
}

func (c *Counter) Increment() (string, error) {
	c.Current += 1
	fmt.Printf("Currently: %d\n", c.Current)
	return fmt.Sprint(c.Current), nil
}

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server on %s: %v", port, err)
	}

	log.Printf("Running on port %s\n", port)

	c := &Counter{Current: 1}

	gs := grpc.NewServer()

	s := gocron_server.NewSchedulerService()
	defer s.Shutdown()

	incr := gocron_server.NewUnitExecFn(c.Increment)
	unit := gocron_server.NewManualUnit("incr", incr)
	err = s.AddUnit(unit)
	if err != nil {
		log.Fatalln(err)
	}

	pb.RegisterSchedulerServer(gs, s)
	reflection.Register(gs)

	if err := gs.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
