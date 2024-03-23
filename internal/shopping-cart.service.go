package internal

import (
	"gorm.io/gorm"
)

type ShoppingCartService struct {
	db *gorm.DB
}

func NewShoppingCartService(db *gorm.DB) *ShoppingCartService {
	return &ShoppingCartService{db}
}

func (s *ShoppingCartService) List() ([]*ShoppingCart, error) {
	shoppingCarts := []ShoppingCart{}
	result := s.db.Find(&shoppingCarts)
	if result.Error != nil {
		return nil, result.Error
	}

	out := []*ShoppingCart{}
	for _, v := range shoppingCarts {
		out = append(out, &v)
	}
	return out, nil
}

func (s *ShoppingCartService) SaveOne(name string) {
	shoppingCart := ShoppingCart{Name: name}
	s.db.Create(&shoppingCart)
}

func (s *ShoppingCartService) FindOne(id uint) ShoppingCart {
	var result ShoppingCart
	s.db.Model(ShoppingCart{Model: gorm.Model{ID: id}}).First(&result)

	return result
}

func (s *ShoppingCartService) DeleteOne(id uint) error {
	return s.db.Delete(&ShoppingCart{}, id).Error
}

func (s *ShoppingCartService) UpdateOne(id uint, payload *ShoppingCart) error {
	return s.db.Model(ShoppingCart{Model: gorm.Model{ID: id}}).Updates(payload).Error
}
