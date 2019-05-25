package data

type BannerInfo interface {
	GetSerial() uint16
	GetEvent() string
	GetText() string
	GetImage() string
	GetURL() string
}
