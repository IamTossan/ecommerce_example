package shoppingcart

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

var (
	CartIsAlreadyClosedError = errors.New("CART_IS_ALREADY_CLOSED")
	UnknownCommandError      = errors.New("UNKNOWN_COMMAND")
)

func OpenShoppingCartHandler(command OpenShoppingCartCommand) (ShoppingCartOpenedEvent, error) {
	return NewShoppingCartOpenedEvent(command.Data.ClientId), nil
}
func CancelShoppingCartHandler(command CancelShoppingCartCommand, aggregate ShoppingCartAggregate) (ShoppingCartCancelledEvent, error) {
	if aggregate.Status != "shopping_cart_opened" {
		return ShoppingCartCancelledEvent{}, CartIsAlreadyClosedError
	}
	return NewShoppingCartCancelled(command.Data.ClientId), nil
}

func HandleCommand(msg *nats.Msg) (Event, error) {
	log.Println(strings.TrimPrefix(msg.Subject, "commands."))
	switch c := strings.TrimPrefix(msg.Subject, "commands."); c {
	case "open_new_shopping_cart":
		var x OpenShoppingCartCommand
		json.Unmarshal(msg.Data, &x)
		event, err := OpenShoppingCartHandler(x)
		if err != nil {
			return nil, err
		}
		return event, nil
	default:
		return nil, UnknownCommandError
	}
}
