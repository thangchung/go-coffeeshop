package repo

import (
	"context"
	"database/sql"

	"github.com/thangchung/go-coffeeshop/internal/barista/domain"
	"github.com/thangchung/go-coffeeshop/internal/barista/infras/postgresql"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"golang.org/x/exp/slog"
)

var _ domain.OrderRepo = (*orderRepo)(nil)

type orderRepo struct {
	pg *postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) domain.OrderRepo {
	return &orderRepo{pg: pg}
}

func (d *orderRepo) Create(ctx context.Context, baristaOrder *domain.BaristaOrder) error {
	slog.Info("create", "domain.BaristaOrder", *baristaOrder)
	// tx, err := d.db.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer tx.Rollback()

	queries := postgresql.New(d.pg.DB)
	// qtx := queries.WithTx(tx)

	slog.Info("debug: itemType", "itemType", baristaOrder.ItemType)
	slog.Info("debug: itemType", "itemType32", int32(baristaOrder.ItemType))

	_, err := queries.CreateOrder(ctx, postgresql.CreateOrderParams{
		ID:       baristaOrder.ID,
		ItemType: int32(baristaOrder.ItemType),
		ItemName: baristaOrder.ItemName,
		TimeUp:   baristaOrder.TimeUp,
		Created:  baristaOrder.Created,
		Updated: sql.NullTime{
			Time:  baristaOrder.Updated,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	return nil

	// tx, err := d.pg.Pool.Begin(ctx)
	// if err != nil {
	// 	return errors.Wrapf(err, "orderRepo-Create-d.pg.Pool.Begin(ctx)")
	// }

	// // insert order
	// sql, args, err := d.pg.Builder.
	// 	Insert(`"barista".barista_orders`).
	// 	Columns("id", "item_type", "item_name", "time_up", "created", "updated").
	// 	Values(
	// 		baristaOrder.ID,
	// 		baristaOrder.ItemType,
	// 		baristaOrder.ItemName,
	// 		baristaOrder.TimeUp,
	// 		baristaOrder.Created,
	// 		baristaOrder.Updated,
	// 	).
	// 	ToSql()
	// if err != nil {
	// 	return tx.Rollback(ctx)
	// }

	// _, err = d.pg.Pool.Exec(ctx, sql, args...)
	// if err != nil {
	// 	return tx.Rollback(ctx)
	// }

	// return tx.Commit(ctx)
}
