package query

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type orderQuery struct {
	repo domain.OrderRepo
}

var _ domain.OrderQuery = (*orderQuery)(nil)

func NewOrderQuery(ctx context.Context, r domain.OrderRepo) domain.OrderQuery {
	return &orderQuery{
		repo: r,
	}
}

func (d *orderQuery) GetAll(ctx context.Context) ([]gen.OrderDto, error) {
	ents, err := d.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewOrderQuery - GetAll - d.repo.GetAll(): %w", err)
	}

	return ents, nil
}
