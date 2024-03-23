package shoppingcart

import (
	"encoding/json"
	"strings"

	"github.com/nats-io/nats.go"
)

var (
// CartIsAlreadyClosedError = errors.New("CART_IS_ALREADY_CLOSED")
// CartNotFoundError        = errors.New("CART_IS_NOT_FOUND")
// UnknownCommandError      = errors.New("UNKNOWN_COMMAND")
)

func ShoppingCartOpenedHandler(event ShoppingCartOpenedEvent) error {

	return nil
}

// func CancelShoppingCartHandler(command CancelShoppingCartCommand, aggregate *ShoppingCartAggregate) (ShoppingCartCancelledEvent, error) {
// 	if aggregate.Status != "shopping_cart_opened" {
// 		return ShoppingCartCancelledEvent{}, CartIsAlreadyClosedError
// 	}
// 	return NewShoppingCartCancelled(command.Data.ClientId), nil
// }

func HandleEvent(msg *nats.Msg, esdbService *EsdbService) error {
	switch c := strings.TrimPrefix(msg.Subject, "commands."); c {
	case "open_new_shopping_cart":
		var x ShoppingCartOpenedEvent
		json.Unmarshal(msg.Data, &x)
		err := ShoppingCartOpenedHandler(x)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}
