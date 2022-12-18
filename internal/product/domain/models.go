package domain

type ItemTypeDto struct {
	Name  string  `json:"name"`
	Type  int     `json:"type"`
	Price float64 `json:"price"`
	Image string  `json:"image"`
}

type ItemDto struct {
	Price float64 `json:"price"`
	Type  int     `json:"type"`
}
