// Package banner contains application service which
// coordinates the use cases for banner entity
package banner

import (
	"context"
	"fmt"
	"time"

	domain "github.com/DzananGanic/banner"
)

// New creates new banner application service
func New(
	bdb domain.BannerDB,
	disp domain.BannerDisplayer,
) *Service {
	return &Service{
		banners: bdb,
		disp:    disp,
	}
}

// Service represents banner application service
type Service struct {
	banners domain.BannerDB
	disp    domain.BannerDisplayer
}

// CreateReq represents create banner request
type CreateReq struct {
	Name                  string
	ScheduledDisplayingAt time.Time
	ExpiresAt             time.Time
}

// Validate validates CreateReq and returns error if the validation fails
func (req *CreateReq) Validate() error {
	if (req.Name == "") || (req.ExpiresAt == time.Time{}) || (req.ScheduledDisplayingAt == time.Time{}) {
		return fmt.Errorf("you must set name, scheduled displaying at, and expires at")
	}
	return nil
}

// CreateResp represents create banner response
type CreateResp struct {
	ID domain.BannerID
}

// Create use case creates a new banner and saves it to the repository
func (s *Service) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	b := domain.Banner{
		Name:                  req.Name,
		ScheduledDisplayingAt: req.ScheduledDisplayingAt,
		ExpiresAt:             req.ExpiresAt,
	}

	id, err := s.banners.Save(b)
	if err != nil {
		return nil, err
	}

	return &CreateResp{
		ID: id,
	}, nil
}

// UpdateReq represents the request to update
// banner properties
type UpdateReq struct {
	ID                    domain.BannerID
	Name                  *string
	ScheduledDisplayingAt *time.Time
	ExpiresAt             *time.Time
}

// Validate validates UpdateReq and returns error if the validation fails
func (req *UpdateReq) Validate() error {
	if req.ID == 0 {
		return fmt.Errorf("you must have banner id")
	}
	return nil
}

// Update use case updates the existing banner and saves it to the repository
func (s *Service) Update(ctx context.Context, req *UpdateReq) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	b, err := s.banners.FetchForID(req.ID)
	if err != nil {
		return err
	}

	// Note: updating fields can be done in a nicer way
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.ScheduledDisplayingAt != nil {
		b.ScheduledDisplayingAt = *req.ScheduledDisplayingAt
	}
	if req.ExpiresAt != nil {
		b.ExpiresAt = *req.ExpiresAt
	}

	_, err = s.banners.Save(*b)

	return err
}

// DisplayResp returns the display banner response
type DisplayResp struct {
	Banner domain.Banner
}

// Display loads available domain banners and finds
// the one that should be shown
func (s *Service) Display(ctx context.Context) (*DisplayResp, error) {
	banner, err := s.disp.DisplayBanner()
	if err != nil {
		return nil, err
	}

	return &DisplayResp{Banner: *banner}, nil
}
