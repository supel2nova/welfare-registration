package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
)

type RefService struct {
	repo *repository.Repo
}

func NewRefService(repo *repository.Repo) *RefService {
	return &RefService{repo: repo}
}

func (s *RefService) Provinces(ctx context.Context) ([]dto.RefItem, error) {
	items, err := s.repo.Provinces(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return items, nil
}

func (s *RefService) Districts(ctx context.Context, provinceCode string) ([]dto.RefItem, error) {
	items, err := s.repo.Districts(ctx, provinceCode)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return items, nil
}

func (s *RefService) Subdistricts(ctx context.Context, districtCode string) ([]dto.RefItem, error) {
	items, err := s.repo.Subdistricts(ctx, districtCode)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return items, nil
}

func (s *RefService) SearchAddress(ctx context.Context, q string) ([]dto.AddressOption, error) {
	q = strings.TrimSpace(q)
	if utf8.RuneCountInString(q) < 2 {
		return []dto.AddressOption{}, nil
	}
	items, err := s.repo.SearchAddress(ctx, q)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return items, nil
}

func (s *RefService) ListUsers(ctx context.Context) ([]dto.DevUser, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return users, nil
}
