package sharedkernel

import (
	"fmt"
)

type OrderSource int8

const (
	OrderSourceCounter OrderSource = iota
	OrderSourceWeb
)

func (e OrderSource) String() string {
	return fmt.Sprintf("%d", int(e))
}

type Status int8

const (
	StatusPlaced Status = iota
	StatusInProcess
	StatusFulfilled
)

func (e Status) String() string {
	return fmt.Sprintf("%d", int(e))
}

type Location int8

const (
	LocationAtlanta Location = iota
	LocationCharlotte
	LocationRaleigh
)

func (e Location) String() string {
	return fmt.Sprintf("%d", int(e))
}

type CommandType int8

const (
	CommandTypePlaceOrder CommandType = iota
)

func (e CommandType) String() string {
	return fmt.Sprintf("%d", int(e))
}

type ItemType int8

const (
	ItemTypeCappuccino ItemType = iota
	ItemTypeCoffeeBlack
	ItemTypeCoffeeWithRoom
	ItemTypeEspresso
	ItemTypeEspressoDouble
	ItemTypeLatte
	ItemTypeCakePop
	ItemTypeCroissant
	ItemTypeMuffin
	ItemTypeCroissantChocolate
)

func (e ItemType) String() string {
	switch e {
	case ItemTypeCappuccino:
		return "CAPPUCCINO"
	case ItemTypeCoffeeBlack:
		return "COFFEE_BLACK"
	case ItemTypeCoffeeWithRoom:
		return "COFFEE_WITH_ROOM"
	case ItemTypeEspresso:
		return "ESPRESSO"
	case ItemTypeEspressoDouble:
		return "ESPRESSO_DOUBLE"
	case ItemTypeLatte:
		return "LATTE"
	case ItemTypeCakePop:
		return "CAKEPOP"
	case ItemTypeCroissant:
		return "CROISSANT"
	case ItemTypeMuffin:
		return "MUFFIN"
	case ItemTypeCroissantChocolate:
		return "CROISSANT_CHOCOLATE"
	default:
		return "CAPPUCCINO"
	}
}
