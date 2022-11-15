package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain"
	"github.com/thangchung/go-coffeeshop/internal/counter/domain/model"
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

func (d *orderRepo) getAllFunc() (string, error) {
	sql, _, err := d.pg.Builder.
		Select(`
			o.id as "o.id", 
			order_source as "o.order_source", 
			loyalty_member_id as "o.loyalty_member_id", 
			order_status as "o.order_status",
			l.id as "l.id",
			item_type as "l.item_type",
			name as "l.name",
			price as "l.price",
			item_status as "l.item_status",
			is_barista_order as "l.is_barista_order",
			o.id as "l.order_id"
			`).
		From(`"order".orders o`).Join(`"order".line_items l ON o.id = l.order_id`).
		Limit(_defaultEntityCap).
		ToSql()

	return sql, err
}

func (d *orderRepo) getByIDFunc(id uuid.UUID) (string, []interface{}, error) {
	return d.pg.Builder.
		Select(`
			o.id as "o.id", 
			order_source as "o.order_source", 
			loyalty_member_id as "o.loyalty_member_id", 
			order_status as "o.order_status",
			l.id as "l.id",
			item_type as "l.item_type",
			name as "l.name",
			price as "l.price",
			item_status as "l.item_status",
			is_barista_order as "l.is_barista_order",
			o.id as "l.order_id"
		`).
		From(`"order".orders o`).Join(`"order".line_items l ON o.id = l.order_id`).
		Where("o.id = ?", id).
		ToSql()
}

func (d *orderRepo) GetAll(ctx context.Context) ([]*domain.Order, error) {
	sql, err := d.getAllFunc()
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Builder: %w", err)
	}

	rows, err := d.pg.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var results []model.OrderListResult
	if err := pgxscan.ScanAll(&results, rows); err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetAll-pgxscan.ScanAll: %w", err)
	}

	uniqueResults := lo.UniqBy(results, func(x model.OrderListResult) string {
		return x.Order.ID.String()
	})
	orders := lo.Map(uniqueResults, func(x model.OrderListResult, _ int) *domain.Order {
		return x.Order
	})
	lineItems := lo.Map(results, func(x model.OrderListResult, _ int) *domain.LineItem {
		return x.LineItem
	})
	entities := make([]*domain.Order, 0, _defaultEntityCap)

	for _, o := range orders {
		order := &domain.Order{
			ID:              o.ID,
			OrderSource:     o.OrderSource,
			LoyaltyMemberID: o.LoyaltyMemberID,
			OrderStatus:     o.OrderStatus,
		}

		filters := lo.Filter(lineItems, func(x *domain.LineItem, _ int) bool {
			return x.OrderID.ID() == o.ID.ID()
		})

		for _, ol := range filters {
			order.LineItems = append(order.LineItems, &domain.LineItem{
				ID:             ol.ID,
				ItemType:       ol.ItemType,
				Name:           ol.Name,
				Price:          ol.Price,
				ItemStatus:     ol.ItemStatus,
				IsBaristaOrder: ol.IsBaristaOrder,
			})
		}

		entities = append(entities, order)
	}

	return entities, nil
}

func (d *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	sql, args, err := d.getByIDFunc(id)
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetByID-r.Builder: %w", err)
	}

	rows, err := d.pg.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetByID-r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var results []model.OrderListResult
	if err := pgxscan.ScanAll(&results, rows); err != nil {
		return nil, fmt.Errorf("NewOrderRepo-GetByID-pgxscan.ScanAll: %w", err)
	}

	uniqueResults := lo.UniqBy(results, func(x model.OrderListResult) string {
		return x.Order.ID.String()
	})

	orders := lo.Map(uniqueResults, func(x model.OrderListResult, _ int) *domain.Order {
		return x.Order
	})
	lineItems := lo.Map(results, func(x model.OrderListResult, _ int) *domain.LineItem {
		return x.LineItem
	})

	if len(orders) == 0 {
		return nil, nil
	}

	order := &domain.Order{
		ID:              orders[0].ID,
		OrderSource:     orders[0].OrderSource,
		LoyaltyMemberID: orders[0].LoyaltyMemberID,
		OrderStatus:     orders[0].OrderStatus,
	}

	for _, ol := range lineItems {
		order.LineItems = append(order.LineItems, &domain.LineItem{
			ID:             ol.ID,
			ItemType:       ol.ItemType,
			Name:           ol.Name,
			Price:          ol.Price,
			ItemStatus:     ol.ItemStatus,
			IsBaristaOrder: ol.IsBaristaOrder,
		})
	}

	return order, nil
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
