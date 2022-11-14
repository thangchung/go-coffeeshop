package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

const _defaultEntityCap = 64

type orderRepo struct {
	pg *postgres.Postgres
}

var _ domain.OrderRepo = (*orderRepo)(nil)

func NewOrderRepo(pg *postgres.Postgres) domain.OrderRepo {
	return &orderRepo{pg: pg}
}

func (d *orderRepo) GetAll(ctx context.Context) ([]gen.OrderDto, error) {
	sql, _, err := d.pg.Builder.
		Select("orders.id").
		From(`"order".orders`).Join(`"order".line_items USING(id)`).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Builder: %w", err)
	}

	rows, err := d.pg.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]gen.OrderDto, 0, _defaultEntityCap)

	for rows.Next() {
		o := gen.OrderDto{}

		err = rows.Scan(&o.Id, &o.OrderSource, &o.LoyaltyMemberId, &o.OrderStatus)
		if err != nil {
			return nil, fmt.Errorf("NewOrderRepo-GetAll-rows.Scan: %w", err)
		}

		entities = append(entities, o)
	}

	return entities, nil
}

func (d *orderRepo) GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	sql, args, err := d.pg.Builder.
		Select("o.id, order_source, loyalty_member_id, order_status").
		From(`"order".orders o`).Join(`"order".line_items l ON o.id = l.order_id`).
		Where("o.id = ?", id).
		GroupBy("o.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Builder: %w", err)
	}

	rows, err := d.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Pool.Query: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0, _defaultEntityCap)

	for rows.Next() {
		o := domain.Order{}

		err = rows.Scan(&o.ID, &o.OrderSource, &o.LoyaltyMemberID, &o.OrderStatus)
		if err != nil {
			return nil, fmt.Errorf("NewOrderRepo-GetAll-rows.Scan: %w", err)
		}

		orders = append(orders, o)
	}

	// continue to load order items
	order := orders[0]
	if len(orders) >= 1 {
		sql, args, err := d.pg.Builder.
			Select("id, item_type, name, price, item_status, is_barista_order").
			From(`"order".line_items`).
			Where("order_id = ?", order.ID).
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Builder: %w", err)
		}

		rows, err = d.pg.Pool.Query(ctx, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Pool.Query: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			o := domain.LineItem{}

			err = rows.Scan(&o.ID, &o.ItemType, &o.Name, &o.Price, &o.ItemStatus, &o.IsBaristaOrder)
			if err != nil {
				return nil, fmt.Errorf("NewOrderRepo-GetAll-rows.Scan: %w", err)
			}

			order.LineItems = append(order.LineItems, o)
		}
	}

	return &order, nil
}

func (d *orderRepo) Create(ctx context.Context, orderModel *gen.OrderDto) error {
	tx, err := d.pg.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrapf(err, "orderRepo-Create-d.pg.Pool.Begin(ctx)")
	}

	// insert order
	sql, args, err := d.pg.Builder.
		Insert(`"order".orders`).
		Columns("id", "order_source", "loyalty_member_id", "order_status", "updated").
		Values(
			orderModel.Id,
			orderModel.OrderSource,
			orderModel.LoyaltyMemberId,
			orderModel.OrderStatus,
			time.Now(),
		).
		ToSql()
	if err != nil {
		return tx.Rollback(ctx)
	}

	_, err = d.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return tx.Rollback(ctx)
	}

	// continue to insert order items
	for _, item := range orderModel.LineItems {
		sql, args, err = d.pg.Builder.
			Insert(`"order".line_items`).
			Columns("id", "item_type", "name", "price", "item_status", "is_barista_order", "order_id", "created", "updated").
			Values(
				uuid.New(),
				item.ItemType,
				item.Name,
				item.Price,
				item.ItemStatus,
				item.IsBaristaOrder,
				orderModel.Id,
				time.Now(),
				time.Now(),
			).
			ToSql()
		if err != nil {
			return tx.Rollback(ctx)
		}

		_, err = d.pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return tx.Rollback(ctx)
		}
	}

	return tx.Commit(ctx)
}

func (d *orderRepo) Update(ctx context.Context, orderModel *gen.OrderDto) (*gen.OrderDto, error) {
	tx, err := d.pg.Pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "orderRepo-Update-d.pg.Pool.Begin(ctx)")
	}

	// update order
	sql, args, err := d.pg.Builder.
		Update(`"order".orders`).
		Set("order_status", orderModel.OrderStatus).
		Set("updated", time.Now()).
		Where("id = ?", orderModel.Id).
		ToSql()
	if err != nil {
		return nil, tx.Rollback(ctx)
	}

	_, err = d.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, tx.Rollback(ctx)
	}

	// continue to update order items
	for _, item := range orderModel.LineItems {
		sql, args, err = d.pg.Builder.
			Update(`"order".line_items`).
			Set("item_status", item.ItemStatus).
			Set("updated", time.Now()).
			Where("id = ?", item.Id).
			ToSql()
		if err != nil {
			return nil, tx.Rollback(ctx)
		}

		_, err = d.pg.Pool.Exec(ctx, sql, args...)
		if err != nil {
			return nil, tx.Rollback(ctx)
		}
	}

	return orderModel, tx.Commit(ctx)
}
