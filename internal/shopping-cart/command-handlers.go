package shoppingcart

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
)

var (
	CartIsAlreadyClosedError = errors.New("CART_IS_ALREADY_CLOSED")
	CartNotFoundError        = errors.New("CART_IS_NOT_FOUND")
	UnknownCommandError      = errors.New("UNKNOWN_COMMAND")
)

func OpenShoppingCartHandler(command OpenShoppingCartCommand) (*ShoppingCartOpenedEvent, error) {
	return NewShoppingCartOpenedEvent(command.StreamId, command.Data.ClientId), nil
}
func CancelShoppingCartHandler(command CancelShoppingCartCommand, aggregate *ShoppingCartAggregate) (ShoppingCartCancelledEvent, error) {
	if aggregate.Status != "shopping_cart_opened" {
		return ShoppingCartCancelledEvent{}, CartIsAlreadyClosedError
	}
	return NewShoppingCartCancelled(command.StreamId), nil
}

func HandleCommand(msg *nats.Msg, esdbService *EsdbService) (Event, error) {
	switch c := strings.TrimPrefix(msg.Subject, "commands."); c {
	case "open_new_shopping_cart":
		var x OpenShoppingCartCommand
		json.Unmarshal(msg.Data, &x)
		event, err := OpenShoppingCartHandler(x)
		if err != nil {
			return nil, err
		}
		return event, nil
	case "cancel_shopping_cart":
		var x CancelShoppingCartCommand
		json.Unmarshal(msg.Data, &x)
		fmt.Println("streamID :", x.StreamId.String())
		aggEvents := esdbService.ReadEvents(x.StreamId.String())
		if len(aggEvents) == 0 {
			return nil, CartNotFoundError
		}
		agg := ShoppingCartFromEvents(aggEvents)
		event, err := CancelShoppingCartHandler(x, agg)
		if err != nil {
			return nil, err
		}
		return event, nil
	default:
		return nil, UnknownCommandError
	}
}
