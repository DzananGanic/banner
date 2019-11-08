package mock

import domain "github.com/DzananGanic/banner"

// BannerDisplayer provides banner displayer repository mock
type BannerDisplayer struct {
	DisplayBannerFn      func() (*domain.Banner, error)
	DisplayBannerInvoked bool
}

// DisplayBanner represents the mock for DisplayBanner banner repository method
func (bdb *BannerDisplayer) DisplayBanner() (*domain.Banner, error) {
	bdb.DisplayBannerInvoked = true
	return bdb.DisplayBannerFn()
}
