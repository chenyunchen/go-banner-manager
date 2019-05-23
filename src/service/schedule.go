package service

import (
	"sync"
	"time"
	"errors"
)

// ---------- ScheduledJobService ---------- //

type fn func(string)

type ScheduledJobService struct {
	jobs     []*ScheduledJob
	handlers map[string]fn
}

func NewScheduledJobService() *ScheduledJobService {
	handlers := make(map[string]fn)
	return &ScheduledJobService{handlers: handlers}
}

func (s *ScheduledJobService) AddJob(route, data string, t time.Time) error {
	handler, ok := s.handlers[route]
	if !ok {
		return errors.New("AddJob|HandlerIsNotExist")
	}

	scheduledJob := NewScheduledJob(data, handler)
	scheduledJob.SetTimer(t)
	s.jobs = append(s.jobs, scheduledJob)

	return nil
}

func (s *ScheduledJobService) Handle(route string, f fn) {
	s.handlers[route] = f
}

func (s *ScheduledJobService) ClearAllJobs() {
	for _, job := range s.jobs {
		job.StopTimer()
	}
	s.jobs = nil

	return
}

// ---------- ScheduledJob ---------- //

type ScheduledJob struct {
	data     string
	handler  fn
	canceled chan bool
	once     sync.Once
}

func NewScheduledJob(data string, f fn) *ScheduledJob {
	return &ScheduledJob{data: data, handler: f}
}

func (s *ScheduledJob) SetTimer(t time.Time) {
	s.StopTimer()
	canceled := make(chan bool)
	go func() {
		timer := time.NewTimer(t.Sub(time.Now()))
		defer timer.Stop()

		select {
		case <-timer.C:
			s.handler(s.data)
		case <-canceled:
			return
		}
	}()

	s.canceled = canceled

	return
}

func (s *ScheduledJob) StopTimer() {
	if s.canceled != nil {
		s.once.Do(func() {
			close(s.canceled)
		})
	}

	return
}