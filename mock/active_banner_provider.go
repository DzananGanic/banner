package mock

import (
	domain "github.com/DzananGanic/banner"
)

// ActiveBannerProvider provides active banner provider repository mock
type ActiveBannerProvider struct {
	SetFn      func(domain.Banner) error
	SetInvoked bool

	GetFn      func() (*domain.Banner, error)
	GetInvoked bool
}

// Set represents set mock implementation
func (a *ActiveBannerProvider) Set(b domain.Banner) error {
	a.SetInvoked = true
	return a.SetFn(b)
}

// Get represents get mock implementation
func (a *ActiveBannerProvider) Get() (*domain.Banner, error) {
	a.GetInvoked = true
	return a.GetFn()
}
