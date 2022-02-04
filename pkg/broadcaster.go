package gocron_server

import (
	pb "github.com/strnophix/gocron-server/pkg/proto"
)

type Streamer interface {
	Send(*pb.ListenJobResponse) error
}

type Streams []Streamer

type EventBroadcaster struct {
	streams Streams
}

func NewEventBroadcaster() *EventBroadcaster {
	return &EventBroadcaster{
		streams: make(Streams, 0),
	}
}

func (b *EventBroadcaster) Subscribe(stream Streamer) {
	b.streams = append(b.streams, stream)
}

func (b *EventBroadcaster) Publish(resp *pb.ListenJobResponse) {
	closed := make([]int, 0)
	for idx, stream := range b.streams {
		if err := stream.Send(resp); err != nil {
			closed = append(closed, idx)
		}
	}

	if len(closed) == len(b.streams) {
		b.streams = nil
		return
	}

	for idx := len(closed) - 1; idx > -1; idx-- {
		removeStream(b.streams, idx)
	}
}

func (b *EventBroadcaster) SubscriberCount() int {
	return len(b.streams)
}

func removeStream(s Streams, i int) Streams {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func NewBroadcastResponse(name, result string) *pb.ListenJobResponse {
	return &pb.ListenJobResponse{JobName: name, JobResult: result}
}
