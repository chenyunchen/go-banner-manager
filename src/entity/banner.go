package entity

// Banner Command
const (
	GetBannersRequest_CMD              = 0x0001
	UpdateBannerRequest_CMD            = 0x0002
	UpdateBannerStartedTimeRequest_CMD = 0x0003
	UpdateBannerExpiredTimeRequest_CMD = 0x0004
	ClearAllBannerTimersRequest_CMD    = 0x0005
)

// UpdateBannerRequest is the structure for update banner request
type UpdateBannerRequest struct {
	Serial      uint16 `json:"serial"`
	StartedTime uint32 `json:"startedTime"`
	ExpiredTime uint32 `json:"expiredTime"`
}

// UpdateBannerStartedTimeRequest is the structure for update banner started time request
type UpdateBannerStartedTimeRequest struct {
	Serial      uint16 `json:"serial"`
	StartedTime uint32 `json:"startedTime"`
}

// UpdateBannerExpiredTimeRequest is the structure for update banner expired time request
type UpdateBannerExpiredTimeRequest struct {
	Serial      uint16 `json:"serial"`
	ExpiredTime uint32 `json:"expiredTime"`
}

// Banner is the structure for basic banner info
type Banner struct {
	Serial      uint16 `json:"serial"`
	Event       string `json:"event"`
	Text        string `json:"text"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	StartedTime string `json:"startedTime"`
	ExpiredTime string `json:"expiredTime"`
}

// GetSerial implements the Banner interface
func (b *Banner) GetSerial() uint16 {
	return b.Serial
}

// GetEvent implements the Banner interface
func (b *Banner) GetEvent() string {
	return b.Event
}

// GetText implements the Banner interface
func (b *Banner) GetText() string {
	return b.Text
}

// GetImage implements the Banner interface
func (b *Banner) GetImage() string {
	return b.Image
}

// GetURL implements the Banner interface
func (b *Banner) GetURL() string {
	return b.URL
}
