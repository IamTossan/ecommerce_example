package shoppingcart

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	isEvent()
	GetType() string
	GetStreamId() string
}

func (e ShoppingCartOpenedEvent) isEvent()             {}
func (e ProductAddedToShoppingCartEvent) isEvent()     {}
func (e ProductRemovedFromShoppingCartEvent) isEvent() {}
func (e ShoppingCartCancelledEvent) isEvent()          {}
func (e ShoppingCartConfirmedEvent) isEvent()          {}

func (e ShoppingCartOpenedEvent) GetType() string {
	return e.Type
}
func (e ProductAddedToShoppingCartEvent) GetType() string {
	return e.Type
}
func (e ProductRemovedFromShoppingCartEvent) GetType() string {
	return e.Type
}
func (e ShoppingCartCancelledEvent) GetType() string {
	return e.Type
}
func (e ShoppingCartConfirmedEvent) GetType() string {
	return e.Type
}
func (e ShoppingCartOpenedEvent) GetStreamId() string {
	return e.StreamId.String()
}
func (e ProductAddedToShoppingCartEvent) GetStreamId() string {
	return e.StreamId.String()
}
func (e ProductRemovedFromShoppingCartEvent) GetStreamId() string {
	return e.StreamId.String()
}
func (e ShoppingCartCancelledEvent) GetStreamId() string {
	return e.StreamId.String()
}
func (e ShoppingCartConfirmedEvent) GetStreamId() string {
	return e.StreamId.String()
}

type EventHeader struct {
	Id             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	Timestamp      time.Time `json:"timestamp"`
	StreamId       uuid.UUID `json:"stream_id"`
	StreamPosition uint      `json:"stream_position"`
}

func NewEventHeader(eventType string) EventHeader {
	eventId, _ := uuid.NewUUID()
	now := time.Now()
	return EventHeader{
		Id:        eventId,
		Type:      eventType,
		Timestamp: now,
	}
}

type ShoppingCartOpenedBody struct {
	ShoppingCartId uuid.UUID `json:"shopping_cart_id"`
	ClientId       uuid.UUID `json:"client_id"`
	OpenedAt       time.Time `json:"opened_at"`
}

type ShoppingCartOpenedEvent struct {
	EventHeader
	Data ShoppingCartOpenedBody `json:"data"`
}

func NewShoppingCartOpenedEvent(streamId, clientId uuid.UUID) *ShoppingCartOpenedEvent {
	eventHeader := NewEventHeader("shopping_cart_opened")
	eventHeader.StreamId = streamId
	eventHeader.StreamPosition = 0

	return &ShoppingCartOpenedEvent{
		EventHeader: eventHeader,
		Data: ShoppingCartOpenedBody{
			ShoppingCartId: streamId,
			ClientId:       clientId,
			OpenedAt:       eventHeader.Timestamp,
		},
	}
}

type ProductItem struct {
	ProductId string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}

type PricedProductItem struct {
	ProductItem ProductItem
	UnitPrice   uint `json:"unit_price"`
}

type ProductWithShoppingCartId struct {
	ShoppingCartId uuid.UUID `json:"shopping_cart_id"`
	ProductItem    PricedProductItem
}

type ProductAddedToShoppingCartEvent struct {
	EventHeader
	Data ProductWithShoppingCartId `json:"data"`
}

func NewProductAddedToShoppingCart(streamId uuid.UUID, pricedProductItem PricedProductItem) ProductAddedToShoppingCartEvent {
	eventHeader := NewEventHeader("product_added_to_shopping_cart")
	eventHeader.StreamId = streamId
	eventHeader.StreamPosition = 0

	return ProductAddedToShoppingCartEvent{
		EventHeader: eventHeader,
		Data: ProductWithShoppingCartId{
			ShoppingCartId: streamId,
			ProductItem:    pricedProductItem,
		},
	}
}

type ProductRemovedFromShoppingCartEvent struct {
	EventHeader
	Data ProductWithShoppingCartId `json:"data"`
}

func NewProductRemovedShoppingCart(streamId uuid.UUID, pricedProductItem PricedProductItem) ProductRemovedFromShoppingCartEvent {
	eventHeader := NewEventHeader("product_removed_from_shopping_cart")
	eventHeader.StreamId = streamId
	eventHeader.StreamPosition = 0

	return ProductRemovedFromShoppingCartEvent{
		EventHeader: eventHeader,
		Data: ProductWithShoppingCartId{
			ShoppingCartId: streamId,
			ProductItem:    pricedProductItem,
		},
	}
}

type ShoppingCartConfirmedBody struct {
	ShoppingCartId uuid.UUID `json:"shopping_cart_id"`
	ConfirmedAt    time.Time `json:"confirmed_at"`
}

type ShoppingCartConfirmedEvent struct {
	EventHeader
	Data ShoppingCartConfirmedBody `json:"data"`
}

func NewShoppingCartConfirmed(streamId uuid.UUID) ShoppingCartConfirmedEvent {
	eventHeader := NewEventHeader("shopping_cart_confirmed")
	eventHeader.StreamId = streamId
	eventHeader.StreamPosition = 0

	return ShoppingCartConfirmedEvent{
		EventHeader: eventHeader,
		Data: ShoppingCartConfirmedBody{
			ShoppingCartId: streamId,
			ConfirmedAt:    time.Now(),
		},
	}
}

type ShoppingCartCancelledBody struct {
	ShoppingCartId uuid.UUID `json:"shopping_cart_id"`
	CancelledAt    time.Time `json:"cancelled_at"`
}

type ShoppingCartCancelledEvent struct {
	EventHeader
	Data ShoppingCartCancelledBody `json:"data"`
}

func NewShoppingCartCancelled(streamId uuid.UUID) ShoppingCartCancelledEvent {
	eventHeader := NewEventHeader("shopping_cart_cancelled")
	eventHeader.StreamId = streamId
	eventHeader.StreamPosition = 0

	return ShoppingCartCancelledEvent{
		EventHeader: eventHeader,
		Data: ShoppingCartCancelledBody{
			ShoppingCartId: streamId,
			CancelledAt:    time.Now(),
		},
	}
}
