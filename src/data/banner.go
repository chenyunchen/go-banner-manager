package data

// BannerInfo is the interface to get basic banner data
type BannerInfo interface {
	GetSerial() uint16
	GetEvent() string
	GetText() string
	GetImage() string
	GetURL() string
}
