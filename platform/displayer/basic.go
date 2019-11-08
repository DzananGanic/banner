package displayer

import (
	"fmt"
	"sort"
	"time"

	domain "github.com/DzananGanic/banner"
)

// NewBasic is factory method that creates new
// banner displayer with basic banner selection algorithm
func NewBasic(
	banners domain.BannerDB,
	activeProvider domain.ActiveBannerProvider,
	ip func() (string, error),
) *BasicBannerDisplayer {
	return &BasicBannerDisplayer{
		banners:        banners,
		activeProvider: activeProvider,
		ip:             ip,
	}
}

// BasicBannerDisplayer represents the basic
// banner provider interface implementation
// It implements basic banner selection algorithm
type BasicBannerDisplayer struct {
	banners        domain.BannerDB
	activeProvider domain.ActiveBannerProvider
	ip             func() (string, error)
}

// DisplayBanner returns the banner that should be shown
func (bp *BasicBannerDisplayer) DisplayBanner() (*domain.Banner, error) {
	abanner, err := bp.activeProvider.Get()
	if err != nil {
		return nil, err
	}

	if !abanner.IsExpired(time.Now()) {
		return abanner, nil
	}

	nextBanner, err := bp.findNextBanner()
	if err != nil {
		return nil, err
	}

	err = bp.activeProvider.Set(*nextBanner)
	if err != nil {
		return nil, err
	}

	return nextBanner, nil
}

func (bp *BasicBannerDisplayer) findNextBanner() (*domain.Banner, error) {
	banners, err := bp.banners.List()
	if err != nil {
		return nil, err
	}

	// we sort the slice by expiration date because of the following requirement:
	// "there may be occasions where two banners are considered active. In this case,
	// the banner with the earlier expiration should be displayed."
	sort.Slice(banners, func(i, j int) bool {
		return banners[i].ExpiresAt.Before(banners[j].ExpiresAt)
	})

	for _, b := range banners {
		// if the banner is expired, just continue
		// in the future we would have a better way of handling this
		// either through database property to filter out expired
		// banners through query or something like that
		if b.IsExpired(time.Now()) {
			continue
		}

		iip, err := bp.ip()
		if err != nil {
			return nil, err
		}

		// if internal ip address is 10.0.0.1 or 10.0.0.2,
		// then it does not matter whether banner is in display period
		if iip == "10.0.0.1" || iip == "10.0.0.2" {
			return &b, nil
		}

		// if not, we just check whether the banner is in display period
		if b.IsInDisplayPeriod(time.Now()) {
			return &b, nil
		}
	}

	return nil, fmt.Errorf("no active banners found")
}
