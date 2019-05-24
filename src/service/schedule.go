package service

import (
	"errors"
	"sync"
	"time"
)

type fn func(string)

//  ScheduledJob represents a job service info
type ScheduledJobService struct {
	jobs     []*ScheduledJob
	handlers map[string]fn
}

// NewRouter creates a new scheduled job service
func NewScheduledJobService() *ScheduledJobService {
	handlers := make(map[string]fn)
	return &ScheduledJobService{handlers: handlers}
}

// AddJob creates a new job to the service
func (s *ScheduledJobService) AddJob(route, data string, t time.Time) error {
	handler, ok := s.handlers[route]
	if !ok {
		return errors.New("ScheduledJobService|AddJob|HandlerIsNotExist")
	}

	scheduledJob := NewScheduledJob(data, handler)
	scheduledJob.SetTimer(t)
	s.jobs = append(s.jobs, scheduledJob)

	return nil
}

// Handle register the handle function for the route
func (s *ScheduledJobService) Handle(route string, f fn) {
	s.handlers[route] = f
}

// ClearAllJobs stop all timer job in the service
func (s *ScheduledJobService) ClearAllJobs() {
	for _, job := range s.jobs {
		job.StopTimer()
	}
	s.jobs = nil

	return
}

// ScheduledJob represents a single job info
type ScheduledJob struct {
	data     string
	handler  fn
	canceled chan bool
	once     sync.Once
}

// NewRouter creates a new scheduled job
func NewScheduledJob(data string, f fn) *ScheduledJob {
	return &ScheduledJob{data: data, handler: f}
}

// SetTimer creates a timer for the job
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

// SetTimer stops a timer for the job
func (s *ScheduledJob) StopTimer() {
	if s.canceled != nil {
		s.once.Do(func() {
			close(s.canceled)
		})
	}

	return
}
