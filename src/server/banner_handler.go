package main

import (
	"encoding/json"
	"log"
	"time"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
)

const (
	Max_Timestamp = 4102358400 // Dec 31 2099. (UTC)
)

// getBannersHandler get banners and if the connection from white list, display inactive banner
func getBannersHandler(ctx *context) {
	service, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	// if there have two active banner, only display the expired time early one.
	banners := make([]entity.Banner, 1)
	minExpiredTime := time.Unix(Max_Timestamp, 0)
	for _, banner := range service.ActiveBanners {
		serial := banner.GetSerial()
		startedTime, expiredTime := service.Schedule.GetJobPeriods(serial)

		if expiredTime.Unix() < minExpiredTime.Unix() {
			minExpiredTime = expiredTime
			banners[0] = entity.Banner{
				Serial:      serial,
				Event:       banner.GetEvent(),
				Text:        banner.GetText(),
				Image:       banner.GetImage(),
				URL:         banner.GetURL(),
				StartedTime: startedTime.In(time.Local).String(),
				ExpiredTime: expiredTime.In(time.Local).String(),
			}
		}
	}

	debug := false
	for _, addr := range service.Config.WhiteList {
		if addr == conn.RemoteAddr() {
			debug = true
			break
		}
	}

	// if there have two active banner, only display the expired time early one.
	if debug {
		for _, serial := range service.Schedule.GetAllInActiveJobSerials() {
			startedTime, expiredTime := service.Schedule.GetJobPeriods(serial)

			if expiredTime.Unix() < minExpiredTime.Unix() {
				minExpiredTime = expiredTime
				banner, err := service.DataManager.GetBanner(serial)
				if err != nil {
					log.Printf("getBannersHandler|service|GetBanner|error:%v", err)
					continue
				}

				banners[0] = entity.Banner{
					Serial:      serial,
					Event:       banner.GetEvent(),
					Text:        banner.GetText(),
					Image:       banner.GetImage(),
					URL:         banner.GetURL(),
					StartedTime: startedTime.In(time.Local).String(),
					ExpiredTime: expiredTime.In(time.Local).String(),
				}
			}
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

// updateBannerHandler update banner with started time and expired time
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
	} else {
		conn.Write([]byte("Timestamp periods is overlap! Only one banner can be actived."))
		return
	}

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

	return
}

// updateBannerStartedTimeHandler update banner started time and expired time is 2099/12/31 by default
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
	} else {
		conn.Write([]byte("Timestamp periods is overlap! Only one banner can be actived."))
		return
	}

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

	return
}

// updateBannerExpiredTimeHandler update banner expired time if started time was set before
func updateBannerExpiredTimeHandler(ctx *context) {
	service, conn, content := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	updateBannerExpiredTimeRequest := entity.UpdateBannerExpiredTimeRequest{}
	err := json.Unmarshal(content, &updateBannerExpiredTimeRequest)
	if err != nil {
		log.Printf("Server|updateBannerExpiredTimeHandler|JsonUnmarshal|error:%v", err)
	}

	startedTime, _ := service.Schedule.GetJobPeriods(updateBannerExpiredTimeRequest.Serial)
	if startedTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		conn.Write([]byte("Start time is not set yet!"))
		return
	}

	timestamp := time.Unix(int64(updateBannerExpiredTimeRequest.ExpiredTime), 0)

	debug := false
	for _, addr := range service.Config.WhiteList {
		if addr == conn.RemoteAddr() {
			debug = true
			break
		}
	}

	isOverlap := service.Schedule.CheckJobPeriodOverlap(debug, updateBannerExpiredTimeRequest.Serial, "expire", timestamp)
	if !isOverlap {
		service.Schedule.AddJob("HideBanner", updateBannerExpiredTimeRequest.Serial, "expire", timestamp)
	} else {
		conn.Write([]byte("Timestamp periods is overlap! Only one banner can be actived."))
		return
	}

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

	return
}

// clearAllBannerTimersHandler clear all the banner timer and all the display banner
func clearAllBannerTimersHandler(ctx *context) {
	service, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	service.Schedule.ClearAllJobs()

	// if the display banner update immediately, delay 200ms to get the result
	time.Sleep(200 * time.Millisecond)
	getBannersHandler(ctx)

	return
}
