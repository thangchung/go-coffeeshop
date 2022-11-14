package command

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type orderCommand struct {
	repo domain.OrderRepo
}

var _ domain.OrderCommand = (*orderCommand)(nil)

func NewOrderCommand(ctx context.Context, r domain.OrderRepo) domain.OrderCommand {
	return &orderCommand{
		repo: r,
	}
}

func (d *orderCommand) Create(ctx context.Context, orderModel *gen.OrderDto) error {
	err := d.repo.Create(ctx, orderModel)
	if err != nil {
		return fmt.Errorf("NewOrderCommand-Create-d.repo.Create(ctx, orderModel): %w", err)
	}

	return nil
}

func (d *orderCommand) Update(ctx context.Context, orderModel *gen.OrderDto) (*gen.OrderDto, error) {
	order, err := d.repo.Update(ctx, orderModel)
	if err != nil {
		return nil, fmt.Errorf("NewOrderCommand-Update-d.repo.Update(ctx, orderModel): %w", err)
	}

	return order, nil
}
