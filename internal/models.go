package internal

import "gorm.io/gorm"

type ShoppingCart struct {
	gorm.Model
	Name string `json:"name"`
}
