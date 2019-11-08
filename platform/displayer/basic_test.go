package displayer_test

import (
	"fmt"
	"testing"
	"time"

	domain "github.com/DzananGanic/banner"
	"github.com/DzananGanic/banner/platform/ip"

	"github.com/DzananGanic/banner/mock"
	"github.com/DzananGanic/banner/platform/displayer"
	"github.com/stretchr/testify/assert"
)

func TestBasicDisplayBanner(t *testing.T) {
	cases := []struct {
		name       string
		now        func() time.Time
		bdb        func() *mock.BannerDB
		ap         func() *mock.ActiveBannerProvider
		ipProvider func() (string, error)
		wantBanner *domain.Banner
		wantErr    bool
	}{
		{
			name: "test active banner provider error",
			bdb: func() *mock.BannerDB {
				return nil
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return nil, fmt.Errorf("database error")
				}
				return active
			},
			ipProvider: ip.Internal,
			wantErr:    true,
		},
		{
			name: "test return currently active banner",
			bdb: func() *mock.BannerDB {
				return nil
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2025, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				return active
			},
			wantBanner: &domain.Banner{
				ExpiresAt: time.Date(2025, 1, 1, 1, 1, 1, 1, time.Local),
			},
			ipProvider: ip.Internal,
			wantErr:    false,
		},
		{
			name: "test active is expired and find next banner throws error",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return nil, fmt.Errorf("database error")
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				return active
			},
			ipProvider: ip.Internal,
			wantErr:    true,
		},
		{
			name: "test active provider set throws error",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return fmt.Errorf("failed setting banner")
				}
				return active
			},
			ipProvider: ip.Internal,
			wantErr:    true,
		},
		{
			name: "test successfully return new banner",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return nil
				}
				return active
			},
			ipProvider: ip.Internal,
			wantErr:    false,
			wantBanner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
			},
		},
		{
			name: "test two active banners, should show one with earlier expiration date",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
						},
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return nil
				}
				return active
			},
			wantErr:    false,
			ipProvider: ip.Internal,
			wantBanner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			},
		},
		{
			name: "test skip unactive banner",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local),
						},
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return nil
				}
				return active
			},
			wantErr:    false,
			ipProvider: ip.Internal,
			wantBanner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			},
		},
		{
			name: "test no active banners found",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local),
						},
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return nil
				}
				return active
			},
			wantErr:    true,
			ipProvider: ip.Internal,
		},
		{
			name: "test error getting ip address",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local),
						},
						{
							ScheduledDisplayingAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2022, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return nil
				}
				return active
			},
			ipProvider: func() (string, error) {
				return "", fmt.Errorf("error getting ip address")
			},
			wantErr: true,
		},
		{
			name: "test show banner if internal IP is 10.0.0.1 even before display period",
			bdb: func() *mock.BannerDB {
				db := &mock.BannerDB{}
				db.ListFn = func() ([]domain.Banner, error) {
					return []domain.Banner{
						{
							ScheduledDisplayingAt: time.Date(2017, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local),
						},
						{
							ScheduledDisplayingAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
							ExpiresAt:             time.Date(2022, 1, 1, 1, 1, 1, 1, time.Local),
						},
					}, nil
				}
				return db
			},
			ap: func() *mock.ActiveBannerProvider {
				active := &mock.ActiveBannerProvider{}
				active.GetFn = func() (*domain.Banner, error) {
					return &domain.Banner{ExpiresAt: time.Date(2008, 1, 1, 1, 1, 1, 1, time.Local)}, nil
				}
				active.SetFn = func(b domain.Banner) error {
					return nil
				}
				return active
			},
			ipProvider: func() (string, error) {
				return "10.0.0.1", nil
			},
			wantErr: false,
			wantBanner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2022, 1, 1, 1, 1, 1, 1, time.Local),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := displayer.NewBasic(
				c.bdb(),
				c.ap(),
				c.ipProvider,
			)

			resp, err := svc.DisplayBanner()
			if c.wantBanner != nil {
				assert.Equal(t, resp, c.wantBanner)
			}
			if c.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
