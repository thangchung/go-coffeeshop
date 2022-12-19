package domain

type OrderListResult struct {
	Order    *Order    `db:"o"`
	LineItem *LineItem `db:"l"`
}
