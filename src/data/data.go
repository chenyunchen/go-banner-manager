package data

type Manager interface {
	BannersManager
}

type BannersManager interface {
	GetBanners() ([]BannerInfo, error)
	GetBanner(uint16) (BannerInfo, error)
}
