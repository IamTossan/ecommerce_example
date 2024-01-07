package shoppingcart

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingCartAggregate struct {
	Id           uuid.UUID           `json:"shopping_cart_id"`
	ClientId     uuid.UUID           `json:"client_id"`
	Status       string              `json:"status"`
	Version      uint                `json:"version"`
	ProductItems []PricedProductItem `json:"product_items"`
	OpenedAt     time.Time           `json:"opened_at"`
	ConfirmedAt  time.Time           `json:"confirmed_at,omitempty"`
	CancelledAt  time.Time           `json:"cancelled_at,omitempty"`
}

func ShoppingCartFromEvents(events []Event) *ShoppingCartAggregate {
	projection := &ShoppingCartAggregate{}

	for _, e := range events {
		projection.On(e)
	}

	return projection
}

func (s *ShoppingCartAggregate) On(event Event) {
	switch e := event.(type) {
	case ShoppingCartOpenedEvent:
		s.Id = e.Data.ShoppingCartId
		s.ClientId = e.Data.ClientId
		s.Status = "shopping_cart_opened"
		s.Version = e.StreamPosition
		s.ProductItems = []PricedProductItem{}
		s.OpenedAt = e.Data.OpenedAt
	case ProductAddedToShoppingCartEvent:
		s.Version = e.StreamPosition
	case ProductRemovedFromShoppingCartEvent:
		s.Version = e.StreamPosition
	case ShoppingCartCancelledEvent:
		s.Status = "shopping_cart_cancelled"
		s.Version = e.StreamPosition
		s.CancelledAt = e.Data.CancelledAt
	case ShoppingCartConfirmedEvent:
		s.Status = "shopping_cart_confirmed"
		s.Version = e.StreamPosition
		s.ConfirmedAt = e.Data.ConfirmedAt
	}
}
