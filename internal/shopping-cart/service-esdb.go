package shoppingcart

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

type EsdbService struct {
	Esdb *esdb.Client
}

func (e *EsdbService) WriteEvent(streamId string, event Event) {
	data, _ := json.Marshal(event)
	eventData := esdb.EventData{
		ContentType: esdb.JsonContentType,
		EventType:   event.GetType(),
		Data:        data,
	}
	_, err := e.Esdb.AppendToStream(context.Background(), streamId, esdb.AppendToStreamOptions{}, eventData)
	if err != nil {
		panic(err)
	}
}

func (e *EsdbService) ReadEvents(streamId string) []Event {
	stream, err := e.Esdb.ReadStream(context.Background(), streamId, esdb.ReadStreamOptions{}, 3)

	if err != nil {
		panic(err)
	}
	defer stream.Close()

	events := []Event{}
	for {
		e, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}
		switch e.Event.EventType {
		case "shopping_cart_opened":
			var x ShoppingCartOpenedEvent
			json.Unmarshal(e.OriginalEvent().Data, &x)
			events = append(events, x)
		case "shopping_cart_cancelled":
			var x ShoppingCartCancelledEvent
			json.Unmarshal(e.OriginalEvent().Data, &x)
			events = append(events, x)

		}
	}
	return events
}
