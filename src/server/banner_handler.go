package main

import (
	"encoding/json"
	"log"
	"time"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
)

func getBannersHandler(ctx *context) {
	service, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	banners := []entity.Banner{}
	for _, banner := range service.DisplayBanner {
		banners = append(banners, entity.Banner{
			Serial: banner.GetSerial(),
			Event:  banner.GetEvent(),
			Text:   banner.GetText(),
			Image:  banner.GetImage(),
			URL:    banner.GetURL(),
		})
	}

	b, err := json.Marshal(banners)
	if err != nil {
		log.Printf("Server|getBannersHandler|JsonMarshal|error:%v", err)
		return
	}
	conn.Write(b)

	return
}

func updateBannerHandler(ctx *context) {
	service, conn, content := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	updateBannerRequest := entity.UpdateBannerRequest{}
	err := json.Unmarshal(content, &updateBannerRequest)
	if err != nil {
		log.Printf("Server|updateBannerHandler|JsonUnmarshal|error:%v", err)
	}

	timestamp := time.Unix(int64(updateBannerRequest.StartedTime), 0)
	service.Schedule.AddJob("DisplayBanner", updateBannerRequest.Serial, "start", timestamp)
	timestamp = time.Unix(int64(updateBannerRequest.ExpiredTime), 0)
	service.Schedule.AddJob("HideBanner", updateBannerRequest.Serial, "expire", timestamp)

	return
}

func updateBannerStartedTimeHandler(ctx *context) {
	service, conn, content := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	updateBannerStartedTimeRequest := entity.UpdateBannerStartedTimeRequest{}
	err := json.Unmarshal(content, &updateBannerStartedTimeRequest)
	if err != nil {
		log.Printf("Server|updateBannerStartedTimeHandler|JsonUnmarshal|error:%v", err)
	}

	timestamp := time.Unix(int64(updateBannerStartedTimeRequest.StartedTime), 0)
	service.Schedule.AddJob("DisplayBanner", updateBannerStartedTimeRequest.Serial, "start", timestamp)

	return
}

func updateBannerExpiredTimeHandler(ctx *context) {
	service, conn, content := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	updateBannerExpiredTimeRequest := entity.UpdateBannerExpiredTimeRequest{}
	err := json.Unmarshal(content, &updateBannerExpiredTimeRequest)
	if err != nil {
		log.Printf("Server|updateBannerExpiredTimeHandler|JsonUnmarshal|error:%v", err)
	}

	timestamp := time.Unix(int64(updateBannerExpiredTimeRequest.ExpiredTime), 0)
	service.Schedule.AddJob("HideBanner", updateBannerExpiredTimeRequest.Serial, "expire", timestamp)

	return
}

func clearAllBannerTimersHandler(ctx *context) {
	service, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	service.Schedule.ClearAllJobs()

	return
}
