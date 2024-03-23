package shoppingcart

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsService struct {
	Nc *nats.Conn
}

func (n *NatsService) SendCommand(ctx context.Context, command Command) error {
	log.Println("in SendCommand")

	payload, err := json.Marshal(command)

	if err != nil {
		return err
	}

	n.Nc.Publish("commands."+command.GetType(), payload)
	return nil
}

func (n *NatsService) SendEvent(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	n.Nc.Publish("events."+event.GetStreamId()+"."+event.GetType(), payload)
	return nil
}
