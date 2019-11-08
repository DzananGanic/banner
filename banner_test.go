package domain_test

import (
	"testing"
	"time"

	domain "github.com/DzananGanic/banner"
	"github.com/stretchr/testify/assert"
)

func TestIsBannerInDisplayPeriod(t *testing.T) {
	chicagoLocation, err := time.LoadLocation("America/Chicago")
	if err != nil {
		panic("failed to get chicago location")
	}

	newYorkLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("failed to get new york location")
	}

	cases := []struct {
		name   string
		now    func() time.Time
		banner *domain.Banner
		want   bool
	}{
		{
			name: "test banner is in display period",
			now: func() time.Time {
				return time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
			},
			banner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "test banner not yet in display period",
			now: func() time.Time {
				return time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
			},
			banner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				ExpiresAt:             time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			},
			want: false,
		},
		{
			name: "test banner after in display period",
			now: func() time.Time {
				return time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
			},
			banner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
				ExpiresAt:             time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local),
			},
			want: false,
		},
		{
			name: "test banner time zone expired",
			now: func() time.Time {
				return time.Date(2019, 5, 5, 13, 0, 0, 0, chicagoLocation)
			},
			banner: &domain.Banner{
				ScheduledDisplayingAt: time.Date(2018, 1, 1, 0, 0, 0, 0, newYorkLocation),
				ExpiresAt:             time.Date(2019, 5, 5, 13, 0, 0, 0, newYorkLocation),
			},
			want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := c.banner.IsInDisplayPeriod(c.now())
			assert.Equal(t, c.want, res)
		})
	}
}

func TestIsBannerExpired(t *testing.T) {
	chicagoLocation, err := time.LoadLocation("America/Chicago")
	if err != nil {
		panic("failed to get chicago location")
	}

	newYorkLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("failed to get new york location")
	}

	cases := []struct {
		name   string
		now    func() time.Time
		banner *domain.Banner
		want   bool
	}{
		{
			name: "test banner is not expired yet",
			now: func() time.Time {
				return time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
			},
			banner: &domain.Banner{
				ExpiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			},
			want: false,
		},
		{
			name: "test banner is expired",
			now: func() time.Time {
				return time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
			},
			banner: &domain.Banner{
				ExpiresAt: time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "test banner time zone expired",
			now: func() time.Time {
				return time.Date(2019, 5, 5, 13, 0, 0, 0, chicagoLocation)
			},
			banner: &domain.Banner{
				ExpiresAt: time.Date(2019, 5, 5, 13, 0, 0, 0, newYorkLocation),
			},
			want: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := c.banner.IsExpired(c.now())
			assert.Equal(t, c.want, res)
		})
	}
}
