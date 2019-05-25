package main

import (
	"encoding/json"
	"log"
)

func getBannersHandler(ctx *context) {
	service, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	banners, err := service.DataManager.GetBanners()
	if err != nil {
		log.Printf("Server|getBannersHandler|service|GetBanners|error:%v", err)
		return
	}

	b, err := json.Marshal(banners)
	if err != nil {
		log.Printf("Server|getBannersHandler|JsonMarshal|error:%v", err)
	}
	conn.Write(b)

	return
}
