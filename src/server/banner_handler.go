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
		serial := banner.GetSerial()
		startedTime, expiredTime := service.Schedule.GetJobPeriods(serial)
		banners = append(banners, entity.Banner{
			Serial:      serial,
			Event:       banner.GetEvent(),
			Text:        banner.GetText(),
			Image:       banner.GetImage(),
			URL:         banner.GetURL(),
			StartedTime: startedTime,
			ExpiredTime: expiredTime,
		})
	}

	debug := false
	for _, addr := range service.Config.WhiteList {
		if addr == conn.RemoteAddr() {
			debug = true
			break
		}
	}

	if debug {
		for _, serial := range service.Schedule.GetAllInActiveJobSerials() {
			startedTime, expiredTime := service.Schedule.GetJobPeriods(serial)
			banner, err := service.DataManager.GetBanner(serial)
			if err != nil {
				log.Printf("getBannersHandler|service|GetBanner|error:%v", err)
				continue
			}

			banners = append(banners, entity.Banner{
				Serial:      serial,
				Event:       banner.GetEvent(),
				Text:        banner.GetText(),
				Image:       banner.GetImage(),
				URL:         banner.GetURL(),
				StartedTime: startedTime,
				ExpiredTime: expiredTime,
			})
		}
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

	startedTime := time.Unix(int64(updateBannerRequest.StartedTime), 0)
	expiredTime := time.Unix(int64(updateBannerRequest.ExpiredTime), 0)

	debug := false
	for _, addr := range service.Config.WhiteList {
		if addr == conn.RemoteAddr() {
			debug = true
			break
		}
	}

	isOverlap := service.Schedule.CheckJobPeriodsOverlap(debug, updateBannerRequest.Serial, startedTime, expiredTime)
	if !isOverlap {
		service.Schedule.AddJob("DisplayBanner", updateBannerRequest.Serial, "start", startedTime)
		service.Schedule.AddJob("HideBanner", updateBannerRequest.Serial, "expire", expiredTime)
	}

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

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

	debug := false
	for _, addr := range service.Config.WhiteList {
		if addr == conn.RemoteAddr() {
			debug = true
			break
		}
	}

	isOverlap := service.Schedule.CheckJobPeriodOverlap(debug, updateBannerStartedTimeRequest.Serial, "start", timestamp)
	if !isOverlap {
		service.Schedule.AddJob("DisplayBanner", updateBannerStartedTimeRequest.Serial, "start", timestamp)
	}

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

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

	debug := false
	for _, addr := range service.Config.WhiteList {
		if addr == conn.RemoteAddr() {
			debug = true
			break
		}
	}

	isOverlap := service.Schedule.CheckJobPeriodOverlap(debug, updateBannerExpiredTimeRequest.Serial, "start", timestamp)
	if !isOverlap {
		service.Schedule.AddJob("HideBanner", updateBannerExpiredTimeRequest.Serial, "expire", timestamp)
	}

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

	return
}

func clearAllBannerTimersHandler(ctx *context) {
	service, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	service.Schedule.ClearAllJobs()

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

	return
}
