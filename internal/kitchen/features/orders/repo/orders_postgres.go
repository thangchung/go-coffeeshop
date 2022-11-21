package repo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/kitchen/domain"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
)

var _ domain.OrderRepo = (*orderRepo)(nil)

type orderRepo struct {
	pg *postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) domain.OrderRepo {
	return &orderRepo{pg: pg}
}

func (d *orderRepo) Create(ctx context.Context, kitchenOrder *domain.KitchenOrder) error {
	tx, err := d.pg.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrapf(err, "orderRepo-Create-d.pg.Pool.Begin(ctx)")
	}

	// insert order
	sql, args, err := d.pg.Builder.
		Insert(`"kitchen".kitchen_orders`).
		Columns("id", "order_id", "item_type", "item_name", "time_up", "created", "updated").
		Values(
			kitchenOrder.ID,
			kitchenOrder.OrderID,
			kitchenOrder.ItemType,
			kitchenOrder.ItemName,
			kitchenOrder.TimeUp,
			kitchenOrder.Created,
			kitchenOrder.Updated,
		).
		ToSql()
	if err != nil {
		return tx.Rollback(ctx)
	}

	_, err = d.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return tx.Rollback(ctx)
	}

	return tx.Commit(ctx)
}
