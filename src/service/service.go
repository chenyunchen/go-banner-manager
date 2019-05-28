package service

import (
	"log"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/config"
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/data"
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/data/fileData"
)

// Service is the structure for manager services
type Service struct {
	Schedule      *ScheduledJobService
	DataManager   data.Manager
	Config        *config.Config
	ActiveBanners map[uint16]data.BannerInfo
}

// New will create service
func New(path string) *Service {
	config, err := config.Read(path)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	dataManager := fileData.New(config.Data)
	activeBanners := make(map[uint16]data.BannerInfo)

	schedule := NewScheduledJobService()
	schedule.Handle("DisplayBanner", func(serial uint16) {
		banner, err := dataManager.GetBanner(serial)
		if err != nil {
			log.Printf("Service|Schedule|DisplayBanner|GetBanner|error:%v", err)
			return
		}
		activeBanners[serial] = banner
	})
	schedule.Handle("HideBanner", func(serial uint16) {
		delete(activeBanners, serial)
	})

	return &Service{
		Config:        &config,
		DataManager:   dataManager,
		ActiveBanners: activeBanners,
		Schedule:      schedule,
	}
}
