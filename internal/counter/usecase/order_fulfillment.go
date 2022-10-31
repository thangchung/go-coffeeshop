package usecase

import (
	"context"
	"fmt"

	gen "github.com/thangchung/go-coffeeshop/proto/gen"
)

type DefaultQueryOrderFulfillmentUseCase struct {
	repo QueryOrderFulfillmentRepo
}

func NewQueryOrderFulfillmentUseCase(r QueryOrderFulfillmentRepo) *DefaultQueryOrderFulfillmentUseCase {
	return &DefaultQueryOrderFulfillmentUseCase{
		repo: r,
	}
}

func (d DefaultQueryOrderFulfillmentUseCase) GetListOrderFulfillment(ctx context.Context) ([]gen.OrderDto, error) {
	entities, err := d.repo.GetListOrderFulfillment(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewQueryOrderFulfillmentUseCase - GetListOrderFulfillment - s.repo.GetListOrderFulfillment: %w", err)
	}

	return entities, nil
}
