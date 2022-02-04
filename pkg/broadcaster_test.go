package gocron_server_test

import (
	"fmt"
	"testing"

	gocron_server "github.com/strnophix/gocron-server/pkg"
	pb "github.com/strnophix/gocron-server/pkg/proto"
)

type TestStreamer struct {
	CalledOnce bool
}

func NewTestStreamer() *TestStreamer {
	return &TestStreamer{CalledOnce: false}
}

func (ts *TestStreamer) Send(stream *pb.ListenJobResponse) error {
	if ts.CalledOnce == true {
		return fmt.Errorf("Send can only be called once")
	}

	ts.CalledOnce = true
	return nil
}

func TestPublish(t *testing.T) {
	cl := NewTestStreamer()
	eb := gocron_server.NewEventBroadcaster()
	eb.Subscribe(cl)
	msg := gocron_server.NewBroadcastResponse("test", "Content :)")
	eb.Publish(msg)

	if cl.CalledOnce == false {
		t.Fatalf("The Send function should have been called on the stream")
	}

	eb.Publish(msg)
	if eb.SubscriberCount() != 0 {
		t.Fatalf("The second call of Send should have unsubscribed the stream")
	}
}
