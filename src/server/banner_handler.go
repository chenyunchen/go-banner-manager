package main

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
	"encoding/json"
	"log"
)

func getBannersHandler(ctx *context) {
	_, conn, _ := ctx.service, ctx.input.Conn, ctx.input.Content
	defer conn.Close()

	// TODO: Remove fake data and use interface implement to get data
	banner := entity.Banner{
		Serial: 1,
		Event:  "Mercari promotion",
		Text:   "20% off",
		Image:  "https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png",
		URL:    "https://www.mercari.com",
	}
	data := []entity.Banner{banner}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Server|getBannersHandler|JsonMarshal|error:%v", err)
	}

	conn.Write(b)
}
