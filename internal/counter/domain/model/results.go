package model

import "github.com/thangchung/go-coffeeshop/internal/counter/domain"

type OrderListResult struct {
	Order    *domain.Order    `db:"o"`
	LineItem *domain.LineItem `db:"l"`
}
