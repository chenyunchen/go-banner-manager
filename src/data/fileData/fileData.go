package fileData

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
	"encoding/json"
	"log"
	"os"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/data"
)

// FileData is the class to set/get/delete data
type FileData struct {
	path string
}

func New(path string) *FileData {
	return &FileData{
		path: path,
	}
}

func (d *FileData) GetBanners() (banners []data.BannerInfo, err error) {
	file, err := os.Open(d.path)
	if err != nil {
		log.Printf("FileData|GetBanners|FileOpen|error:%v", err)
		return
	}
	defer file.Close()

	b := []entity.Banner{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&b)
	if err != nil {
		log.Printf("FileData|GetBanners|JsonDecode|error:%v", err)
		return
	}

	for _, banner := range b {
		banners = append(banners, &banner)
	}

	return
}
