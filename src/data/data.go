package data

// Manager is the interface to manage all data manager
type Manager interface {
	BannersManager
}

// BannersManager is the interface to manage banner
type BannersManager interface {
	GetBanners() ([]BannerInfo, error)
	GetBanner(uint16) (BannerInfo, error)
}
