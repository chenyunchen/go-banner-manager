package entity

// CMD
const (
	GetBannersRequest_CMD              = 0x0001
	UpdateBannerRequest_CMD            = 0x0002
	UpdateBannerStartedTimeRequest_CMD = 0x0003
	UpdateBannerExpiredTimeRequest_CMD = 0x0004
	ClearAllBannerTimersRequest_CMD    = 0x0005
)

type UpdateBannerRequest struct {
	Serial      uint16 `json:"serial"`
	StartedTime uint32 `json:"startedTime"`
	ExpiredTime uint32 `json:"expiredTime"`
}

type UpdateBannerStartedTimeRequest struct {
	Serial      uint16 `json:"serial"`
	StartedTime uint32 `json:"startedTime"`
}

type UpdateBannerExpiredTimeRequest struct {
	Serial      uint16 `json:"serial"`
	ExpiredTime uint32 `json:"expiredTime"`
}

type Banner struct {
	Serial      uint16 `json:"serial"`
	Event       string `json:"event"`
	Text        string `json:"text"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	StartedTime string `json:"startedTime"`
	ExpiredTime string `json:"expiredTime"`
}

func (b *Banner) GetSerial() uint16 {
	return b.Serial
}

func (b *Banner) GetEvent() string {
	return b.Event
}

func (b *Banner) GetText() string {
	return b.Text
}

func (b *Banner) GetImage() string {
	return b.Image
}

func (b *Banner) GetURL() string {
	return b.URL
}
