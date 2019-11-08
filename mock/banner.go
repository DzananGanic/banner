package mock

import (
	domain "github.com/DzananGanic/banner"
)

// BannerDB provides banner repository mock
type BannerDB struct {
	SaveFn      func(b domain.Banner) (domain.BannerID, error)
	SaveInvoked bool

	FetchForIDFn      func(id domain.BannerID) (*domain.Banner, error)
	FetchForIDInvoked bool

	ListFn      func() ([]domain.Banner, error)
	ListInvoked bool

	DeleteFn      func(domain.BannerID) error
	DeleteInvoked bool
}

// Save represents the mock for Save banner repository method
func (bdb *BannerDB) Save(b domain.Banner) (domain.BannerID, error) {
	bdb.SaveInvoked = true
	return bdb.SaveFn(b)
}

// FetchForID represents the mock for FetchForID banner repository method
func (bdb *BannerDB) FetchForID(id domain.BannerID) (*domain.Banner, error) {
	bdb.FetchForIDInvoked = true
	return bdb.FetchForIDFn(id)
}

// List represents the mock for List banner repository method
func (bdb *BannerDB) List() ([]domain.Banner, error) {
	bdb.ListInvoked = true
	return bdb.ListFn()
}

// Delete represents the mock for Delete banner repository method
func (bdb *BannerDB) Delete(id domain.BannerID) error {
	bdb.DeleteInvoked = true
	return bdb.DeleteFn(id)
}
