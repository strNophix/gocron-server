package gocron_server_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	gocron_server "github.com/strnophix/gocron-server/pkg"
	pb "github.com/strnophix/gocron-server/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var c *Counter

func init() {
	c = &Counter{Current: 1}

	lis = bufconn.Listen(bufSize)
	gs := grpc.NewServer()
	s := gocron_server.NewSchedulerService()
	defer s.Shutdown()

	incr := gocron_server.NewUnitExecFn(c.Increment)
	unit := gocron_server.NewManualUnit("incr", incr)
	s.AddUnit(unit)

	pb.RegisterSchedulerServer(gs, s)
	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestRunJob(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	client := pb.NewSchedulerClient(conn)

	_, err = client.RunJob(ctx, &pb.RunJobRequest{UnitName: "niku"})
	if err == nil {
		t.Fatalf("RunJob should have returned an error for unit niku")
	}

	_, err = client.RunJob(ctx, &pb.RunJobRequest{UnitName: "incr"})
	if err != nil {
		t.Fatalf("RunJob call should have passed but got: %v", err)
	}

	time.Sleep(1 * time.Second)
	if c.Current != 2 {
		t.Fatalf("RunJob call `incr` should have incremented counter but stays at: %d", c.Current)
	}
}
