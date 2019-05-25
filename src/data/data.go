package data

type Manager interface {
	BannersManager
}

type BannersManager interface {
	GetBanners() ([]BannerInfo, error)
}
