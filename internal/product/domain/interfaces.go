package domain

import (
	"context"
)

type (
	ProductRepo interface {
		GetAll(context.Context) ([]*ItemTypeDto, error)
		GetByTypes(context.Context, []string) ([]*ItemDto, error)
	}
)
