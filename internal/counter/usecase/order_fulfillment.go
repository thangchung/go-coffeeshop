package usecase

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type DefaultQueryOrderFulfillmentUseCase struct {
	ctx  context.Context
	repo domain.QueryOrderFulfillmentRepo
}

func NewQueryOrderFulfillmentUseCase(ctx context.Context, r domain.QueryOrderFulfillmentRepo) *DefaultQueryOrderFulfillmentUseCase {
	return &DefaultQueryOrderFulfillmentUseCase{
		ctx:  ctx,
		repo: r,
	}
}

func (d DefaultQueryOrderFulfillmentUseCase) GetListOrderFulfillment() ([]gen.OrderDto, error) {
	entities, err := d.repo.GetListOrderFulfillment()
	if err != nil {
		return nil, fmt.Errorf("NewQueryOrderFulfillmentUseCase - GetListOrderFulfillment - s.repo.GetListOrderFulfillment: %w", err)
	}

	return entities, nil
}
