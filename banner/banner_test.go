package banner_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	domain "github.com/DzananGanic/banner"
	"github.com/DzananGanic/banner/banner"
	"github.com/DzananGanic/banner/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name    string
		req     *banner.CreateReq
		wantID  domain.BannerID
		wantErr bool
	}{
		{
			name: "successfully create",
			req: &banner.CreateReq{
				Name:                  "domain Banner",
				ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			},
			wantID:  domain.BannerID(1),
			wantErr: false,
		},
		{
			name: "failed validation no name",
			req: &banner.CreateReq{
				ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			},
			wantID:  0,
			wantErr: true,
		},
		{
			name: "failed create database error",
			req: &banner.CreateReq{
				Name:                  "Fake Banner",
				ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			},
			wantID:  0,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			args := makeBannerArgs()
			svc := banner.New(
				args.bannerDB,
				args.disp,
			)

			resp, err := svc.Create(context.Background(), c.req)
			if c.wantID != 0 {
				assert.Equal(t, c.wantID, resp.ID)
			}
			if c.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name    string
		req     func() *banner.UpdateReq
		wantErr bool
	}{
		{
			name: "successfully update name",
			req: func() *banner.UpdateReq {
				name := "updated name"
				return &banner.UpdateReq{
					ID:   domain.BannerID(2),
					Name: &name,
				}
			},
			wantErr: false,
		},
		{
			name: "successfully update all",
			req: func() *banner.UpdateReq {
				name := "updated name"
				sch := time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local)
				exp := time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local)
				return &banner.UpdateReq{
					ID:                    domain.BannerID(2),
					ScheduledDisplayingAt: &sch,
					ExpiresAt:             &exp,
					Name:                  &name,
				}
			},
			wantErr: false,
		},
		{
			name: "failed validation no id",
			req: func() *banner.UpdateReq {
				name := "updated name"
				return &banner.UpdateReq{
					Name: &name,
				}
			},
			wantErr: true,
		},
		{
			name: "failed update non existing banner",
			req: func() *banner.UpdateReq {
				name := "updated name"
				return &banner.UpdateReq{
					ID:   domain.BannerID(-5),
					Name: &name,
				}
			},
			wantErr: true,
		},
		{
			name: "failed update database error",
			req: func() *banner.UpdateReq {
				name := "updated name, about to fail"
				return &banner.UpdateReq{
					ID:   domain.BannerID(3),
					Name: &name,
				}
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			args := makeBannerArgs()
			svc := banner.New(
				args.bannerDB,
				args.disp,
			)

			err := svc.Update(context.Background(), c.req())
			if c.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDisplay(t *testing.T) {
	cases := []struct {
		name       string
		wantBanner domain.Banner
		wantErr    bool
	}{
		{
			name: "successfully display banner",
			wantBanner: domain.Banner{
				ID:   1,
				Name: "Best banner",
			},
			wantErr: false,
		},
		// due to the simplicity of this method
		// I believe that we do not have to cover
		// database error test for this use case
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			args := makeBannerArgs()
			svc := banner.New(
				args.bannerDB,
				args.disp,
			)

			resp, err := svc.Display(context.Background())
			assert.Equal(t, resp.Banner, c.wantBanner)
			if c.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

type bannerArgs struct {
	bannerDB *mock.BannerDB
	disp     *mock.BannerDisplayer
}

func makeBannerArgs() bannerArgs {
	bannerDB := &mock.BannerDB{}
	disp := &mock.BannerDisplayer{}

	bannerDB.SaveFn = func(b domain.Banner) (domain.BannerID, error) {
		switch b {
		case domain.Banner{
			Name:                  "domain Banner",
			ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
		}:
			return domain.BannerID(1), nil
		case domain.Banner{
			Name:                  "Fake Banner",
			ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
		}:
			return domain.BannerID(0), fmt.Errorf("banner creation failed")
		case domain.Banner{
			ID:                    domain.BannerID(2),
			Name:                  "updated name",
			CreatedAt:             time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
		}:
			return domain.BannerID(2), nil
		case domain.Banner{
			ID:                    domain.BannerID(2),
			Name:                  "updated name",
			CreatedAt:             time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ExpiresAt:             time.Date(2021, 1, 1, 1, 1, 1, 1, time.Local),
		}:
			return domain.BannerID(2), nil
		case domain.Banner{
			ID:                    domain.BannerID(3),
			Name:                  "updated name, about to fail",
			CreatedAt:             time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
			ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
		}:
			return domain.BannerID(0), fmt.Errorf("database error")
		}

		return 0, fmt.Errorf("no matching cases")
	}

	bannerDB.FetchForIDFn = func(id domain.BannerID) (*domain.Banner, error) {
		switch id {
		case domain.BannerID(2):
			return &domain.Banner{
				ID:                    domain.BannerID(2),
				Name:                  "Deprecated name",
				CreatedAt:             time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			}, nil
		case domain.BannerID(-5):
			return nil, fmt.Errorf("non existing banner")
		case domain.BannerID(3):
			return &domain.Banner{
				ID:                    domain.BannerID(3),
				Name:                  "banner that is doomed to fail",
				CreatedAt:             time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ScheduledDisplayingAt: time.Date(2019, 1, 1, 1, 1, 1, 1, time.Local),
				ExpiresAt:             time.Date(2020, 1, 1, 1, 1, 1, 1, time.Local),
			}, nil
		}

		return nil, fmt.Errorf("no matching cases")
	}

	disp.DisplayBannerFn = func() (*domain.Banner, error) {
		return &domain.Banner{
			ID:   1,
			Name: "Best banner",
		}, nil
	}

	return bannerArgs{
		bannerDB: bannerDB,
		disp:     disp,
	}
}
