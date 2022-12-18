package products

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/product/domain"
)

type service struct {
	repo domain.ProductRepo
}

var _ UseCase = (*service)(nil)

func NewService(repo domain.ProductRepo) UseCase {
	return &service{
		repo: repo,
	}
}

func (s *service) GetItemTypes(ctx context.Context) ([]*domain.ItemTypeDto, error) {
	results, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetItemTypes")
	}

	return results, nil
}

func (s *service) GetItemsByType(ctx context.Context, itemTypes string) ([]*domain.ItemDto, error) {
	types := strings.Split(itemTypes, ",")

	results, err := s.repo.GetByTypes(ctx, types)
	if err != nil {
		return nil, errors.Wrap(err, "service.GetItemsByType")
	}

	return results, nil
}
