package service

import (
	"errors"
	"sync"
	"time"
)

const (
	Min_Timestamp = 28800      // Jan 01 1970. (UTC)
	Max_Timestamp = 4102358400 // Dec 31 2099. (UTC)
)

type fn func(uint16)

//  ScheduledJob represents a job service info
type ScheduledJobService struct {
	jobs     []*ScheduledJob
	handlers map[string]fn
	periods  map[uint16][]time.Time
}

// NewRouter creates a new scheduled job service
func NewScheduledJobService() *ScheduledJobService {
	handlers := make(map[string]fn)
	periods := make(map[uint16][]time.Time)
	return &ScheduledJobService{handlers: handlers, periods: periods}
}

// GetJobPeriods get a job periods
func (s *ScheduledJobService) GetJobPeriods(serial uint16) (startedTime time.Time, expiredTime time.Time) {
	_, ok := s.periods[serial]
	if !ok {
		return
	}

	return s.periods[serial][0], s.periods[serial][1]
}

// GetAllInActiveJobSerials get all inactive job serials
func (s *ScheduledJobService) GetAllInActiveJobSerials() (serials []uint16) {
	for serial, period := range s.periods {
		if period[0].Unix() > time.Now().Unix() {
			serials = append(serials, serial)
		}
	}

	return
}

// CheckJobPeriodExist checks a new job if period exist
func (s *ScheduledJobService) CheckJobPeriodOverlap(debug bool, serial uint16, tag string, t time.Time) bool {
	_, ok := s.periods[serial]
	if !ok {
		s.periods[serial] = make([]time.Time, 2)
		s.periods[serial][0] = time.Unix(Min_Timestamp, 0)
		s.periods[serial][1] = time.Unix(Max_Timestamp, 0)
	}

	if tag == "start" {
		// if started time bigger than expired time
		if t.Unix() >= s.periods[serial][1].Unix() {
			return true
		}
		s.periods[serial][0] = t
	} else {
		// if expired time smaller than started time
		if t.Unix() <= s.periods[serial][0].Unix() {
			return true
		}
		s.periods[serial][1] = t
	}

	// Debug mode can active two banner if timestamp range overlap
	if debug {
		return false
	}

	for key, period := range s.periods {
		if key != serial {
			// Check if timestamp range overlap
			if (period[0].Unix() <= s.periods[serial][1].Unix()) && (s.periods[serial][0].Unix() <= period[1].Unix()) {
				return true
			}
		}
	}

	return false
}

// CheckJobPeriodsExist checks a new job if periods exist
func (s *ScheduledJobService) CheckJobPeriodsOverlap(debug bool, serial uint16, startedTime, expiredTime time.Time) bool {
	_, ok := s.periods[serial]
	if !ok {
		s.periods[serial] = make([]time.Time, 2)
	}

	// if started time bigger than expired time
	if startedTime.Unix() >= expiredTime.Unix() {
		return true
	}
	s.periods[serial][0] = startedTime
	s.periods[serial][1] = expiredTime

	// Debug mode can active two banner if timestamp range overlap
	if debug {
		return false
	}

	for key, period := range s.periods {
		if key != serial {
			// Check if timestamp range overlap
			if (period[0].Unix() <= expiredTime.Unix()) && (startedTime.Unix() <= period[1].Unix()) {
				return true
			}
		}
	}

	return false
}

// AddJob creates a new job to the service
func (s *ScheduledJobService) AddJob(route string, serial uint16, tag string, t time.Time) error {
	handler, ok := s.handlers[route]
	if !ok {
		return errors.New("ScheduledJobService|AddJob|HandlerIsNotExist")
	}

	for _, job := range s.jobs {
		if job.serial == serial && job.tag == tag {
			job.StopTimer()
			job.SetTimer(t)
			return nil
		}
	}

	now := time.Now()
	// Check if started time already begin and the job is not expired yet
	if tag == "start" && t.Unix() <= now.Unix() && s.periods[serial][1].Unix() > now.Unix() {
		t = now
	}
	scheduledJob := NewScheduledJob(serial, tag, handler)
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
	serial   uint16
	tag      string
	handler  fn
	canceled chan bool
	once     sync.Once
}

// NewRouter creates a new scheduled job
func NewScheduledJob(serial uint16, tag string, f fn) *ScheduledJob {
	return &ScheduledJob{serial: serial, tag: tag, handler: f}
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
			s.handler(s.serial)
		case <-canceled:
			return
		}
	}()

	s.canceled = canceled

	return
}

// SetTimer stops a timer for the job
func (s *ScheduledJob) StopTimer() {
	if s.tag == "expire" {
		s.handler(s.serial)
	}
	if s.canceled != nil {
		s.once.Do(func() {
			close(s.canceled)
		})
	}

	return
}
