package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config is the structure for server
type Config struct {
	Data      string   `json:"data"`
	WhiteList []string `json:"whiteList"`
}

// Read will read config file
func Read(path string) (c Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Config|Read|FileOpen|error:%v", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Printf("Config|Read|JsonDecode|error:%v", err)
		return
	}

	return
}
