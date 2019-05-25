package service

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/config"
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/data"
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/data/fileData"
	"log"
)

type Service struct {
	Schedule    *ScheduledJobService
	DataManager data.Manager
	Config      *config.Config
}

func New() *Service {
	config, err := config.Read("./config/local.json")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &Service{
		Schedule:    NewScheduledJobService(),
		DataManager: fileData.New(config.Data),
		Config:      &config,
	}
}
