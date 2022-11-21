package domain

import (
	"context"

	"github.com/thangchung/go-coffeeshop/proto/gen"
)

type (
	ProductRepo interface {
		GetAll(context.Context) ([]*gen.ItemTypeDto, error)
		GetByTypes(context.Context, []string) ([]*gen.ItemDto, error)
	}
)
