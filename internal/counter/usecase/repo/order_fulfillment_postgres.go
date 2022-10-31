package repo

import (
	"context"
	"fmt"

	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

const _defaultEntityCap = 64

type DefaultQueryOrderFulfillmentRepo struct {
	*postgres.Postgres
}

func NewQueryOrderFulfillmentRepo(pg *postgres.Postgres) *DefaultQueryOrderFulfillmentRepo {
	return &DefaultQueryOrderFulfillmentRepo{pg}
}

func (d DefaultQueryOrderFulfillmentRepo) GetListOrderFulfillment(ctx context.Context) ([]gen.OrderDto, error) {
	sql, _, err := d.Builder.
		Select("orders.id").
		From(`"order".orders`).Join(`"order".line_items USING(id)`).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("DefaultQueryOrderFulfillmentRepo - GetListOrderFulfillment - r.Builder: %w", err)
	}

	rows, err := d.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("DefaultQueryOrderFulfillmentRepo - GetListOrderFulfillment - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]gen.OrderDto, 0, _defaultEntityCap)

	for rows.Next() {
		o := gen.OrderDto{}

		err = rows.Scan(&o.Id)
		if err != nil {
			return nil, fmt.Errorf("DefaultQueryOrderFulfillmentRepo - GetListOrderFulfillment - rows.Scan: %w", err)
		}

		entities = append(entities, o)
	}

	return entities, nil
}
