package gocron_server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	pb "github.com/strnophix/gocron-server/pkg/proto"
)

type JobFunc func() (string, error)
type UnitStore map[string]*SchedulerUnit
type JobStore map[string]*gocron.Job

type SchedulerService struct {
	pb.UnimplementedSchedulerServer
	UnitStore
	JobStore
	EventBroadcaster

	Scheduler *gocron.Scheduler
}

func NewSchedulerService() *SchedulerService {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	return &SchedulerService{
		Scheduler: scheduler,
		UnitStore: make(UnitStore),
		JobStore:  make(JobStore),
	}
}

func (s *SchedulerService) Shutdown() {
	s.Scheduler.Stop()
}

func (s *SchedulerService) AddUnit(unit *SchedulerUnit) error {
	s.UnitStore[unit.Name] = unit

	if unit.Cron != "" {
		routine := s.BuildRoutine(unit)
		job, err := s.Scheduler.Cron(unit.Cron).SingletonMode().Do(routine)
		if err != nil {
			return err
		}
		s.JobStore[unit.Name] = job
	}
	return nil
}

func (s *SchedulerService) BuildRoutine(unit *SchedulerUnit) func() {
	return func() {
		out, err := unit.Exec.Call()

		if err != nil {
			msg := NewBroadcastResponse(unit.Name, err.Error())
			s.EventBroadcaster.Publish(msg)
			return
		}

		msg := NewBroadcastResponse(unit.Name, out)
		s.EventBroadcaster.Publish(msg)
	}
}

func NewRunJobError(reason string) (*pb.RunJobResponse, error) {
	return &pb.RunJobResponse{}, fmt.Errorf(reason)
}

func NewRunJobSucces() (*pb.RunJobResponse, error) {
	return &pb.RunJobResponse{}, nil
}

func (s *SchedulerService) RunJob(ctx context.Context, req *pb.RunJobRequest) (*pb.RunJobResponse, error) {
	unit, exists := s.UnitStore[req.UnitName]
	if !exists {
		return NewRunJobError(fmt.Sprintf("Unit with name %s does not exist", req.UnitName))
	}

	routine := s.BuildRoutine(unit)

	if req.RunAt != 0 {
		ts := time.Unix(req.RunAt, 0).UTC()
		job, err := s.Scheduler.Every(1).Day().At(ts).LimitRunsTo(1).SingletonMode().Do(routine)

		if err != nil {
			fmt.Printf("Unix run error: %v", err)
		}

		s.JobStore[unit.Name] = job
		return NewRunJobSucces()
	}

	go routine()

	return NewRunJobSucces()
}

func (s *SchedulerService) ListenJobs(req *pb.ListenJobRequest, stream pb.Scheduler_ListenJobsServer) error {
	s.EventBroadcaster.Subscribe(stream)
	<-stream.Context().Done()
	return nil
}
