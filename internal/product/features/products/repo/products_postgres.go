package repo

import (
	"context"

	"github.com/thangchung/go-coffeeshop/internal/product/domain"
	"github.com/thangchung/go-coffeeshop/proto/gen"
)

var _ domain.ProductRepo = (*productRepo)(nil)

type productRepo struct {
	itemTypes map[string]gen.ItemTypeDto
}

func NewOrderRepo() domain.ProductRepo {
	return &productRepo{
		itemTypes: map[string]gen.ItemTypeDto{
			"CAPPUCCINO": {
				Name:  "CAPPUCCINO",
				Type:  0,
				Price: 4.5,
			},
			"COFFEE_BLACK": {
				Name:  "COFFEE_BLACK",
				Type:  1,
				Price: 3,
			},
			"COFFEE_WITH_ROOM": {
				Name:  "COFFEE_WITH_ROOM",
				Type:  2,
				Price: 3,
			},
			"ESPRESSO": {
				Name:  "ESPRESSO",
				Type:  3,
				Price: 3.5,
			},
			"ESPRESSO_DOUBLE": {
				Name:  "ESPRESSO_DOUBLE",
				Type:  4,
				Price: 4.5,
			},
			"LATTE": {
				Name:  "LATTE",
				Type:  5,
				Price: 4.5,
			},
			"CAKEPOP": {
				Name:  "CAKEPOP",
				Type:  6,
				Price: 2.5,
			},
			"CROISSANT": {
				Name:  "CROISSANT",
				Type:  7,
				Price: 3.25,
			},
			"MUFFIN": {
				Name:  "MUFFIN",
				Type:  8,
				Price: 3,
			},
			"CROISSANT_CHOCOLATE": {
				Name:  "CROISSANT_CHOCOLATE",
				Type:  9,
				Price: 3.5,
			},
		},
	}
}

func (p *productRepo) GetAll(ctx context.Context) ([]*gen.ItemTypeDto, error) {
	results := make([]*gen.ItemTypeDto, 0)

	for _, v := range p.itemTypes {
		results = append(results, &gen.ItemTypeDto{
			Name:  v.Name,
			Type:  v.Type,
			Price: v.Price,
		})
	}

	return results, nil
}

func (p *productRepo) GetByTypes(ctx context.Context, itemTypes []string) ([]*gen.ItemDto, error) {
	results := make([]*gen.ItemDto, 0)

	for _, itemType := range itemTypes {
		item := p.itemTypes[itemType]
		if item.Name != "" {
			results = append(results, &gen.ItemDto{
				Price: item.Price,
				Type:  item.Type,
			})
		}
	}

	return results, nil
}
