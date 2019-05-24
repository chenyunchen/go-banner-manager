package service

type Service struct {
	Schedule *ScheduledJobService
}

func New() *Service {
	return &Service{
		Schedule: NewScheduledJobService(),
	}
}
