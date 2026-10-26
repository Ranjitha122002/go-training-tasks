package model

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}
