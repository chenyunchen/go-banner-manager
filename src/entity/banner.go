package entity

// CMD
const (
	GetBannersRequest_CMD = 0x0001
)

type Banner struct {
	Serial uint16 `json:"serial"`
	Event  string `json:"event"`
	Text   string `json:"text"`
	Image  string `json:"image"`
	URL    string `json:"url"`
}

type UpdateBannerRequest struct {
	Serial      uint16 `json:"serial"`
	StartedTime uint32 `json:"startedTime"`
	ExpiredTime uint32 `json:"expiredTime"`
}
