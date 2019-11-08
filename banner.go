// Package domain holds things like
// domain entities, domain service interfaces
// and repository interfaces
//
// For this use case it holds only banner entity with
// its tests
//
// This package should have no external dependencies
// (except for stdlib, like time etc...), and should not import
// any other package from the project
package domain

import (
	"time"
)

// BannerID represents the Banner identifier
type BannerID int64

// Banner represents the banner entity for this use case
type Banner struct {
	ID                    BannerID
	Name                  string
	CreatedAt             time.Time
	ScheduledDisplayingAt time.Time
	ExpiresAt             time.Time
}

// IsInDisplayPeriod checks whether the banner is in display period
func (b *Banner) IsInDisplayPeriod(now time.Time) bool {
	return now.In(time.Local).After(b.ScheduledDisplayingAt.In(time.Local)) &&
		now.Before(b.ExpiresAt.In(time.Local))
}

// IsExpired checks banner expiration date and returns
// boolean on whether the banner is expired
func (b *Banner) IsExpired(now time.Time) bool {
	return now.In(time.Local).After(b.ExpiresAt.In(time.Local))
}

// BannerDB represents Banner entity repository
type BannerDB interface {
	Save(Banner) (BannerID, error)
	FetchForID(BannerID) (*Banner, error)
	List() ([]Banner, error)
	Delete(BannerID) error
}

// ActiveBannerProvider is the repository which
// sets and gets active banner.
type ActiveBannerProvider interface {
	Set(Banner) error
	Get() (*Banner, error)
}

// BannerDisplayer returns the optimal banner to be shown
type BannerDisplayer interface {
	DisplayBanner() (*Banner, error)
}
