package application

import (
	"context"
	"log"
	"runtime"

	shoppingcart "github.com/IamTossan/ecommerce_example/internal/shopping-cart"
	"github.com/nats-io/nats.go"
)

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func (a *App) StartSubscribeCommands(ctx context.Context, subject string) {
	esdbService := &shoppingcart.EsdbService{
		Esdb: a.es,
	}
	natsService := &shoppingcart.NatsService{
		Nc: a.nc,
	}

	i := 0

	a.nc.Subscribe(subject, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		msg.InProgress()
		event, err := shoppingcart.HandleCommand(msg, esdbService)
		if err != nil {
			log.Fatal(err)
		} else {
			natsService.SendEvent(ctx, event)
			esdbService.WriteEvent(event.GetStreamId(), event)
		}
		msg.Ack()
	})

	a.nc.Flush()

	if err := a.nc.LastError(); err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}

func (a *App) StartSubscribeEvents(ctx context.Context, subject string) {
	esdbService := &shoppingcart.EsdbService{
		Esdb: a.es,
	}

	i := 0

	a.nc.Subscribe(subject, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		msg.InProgress()
		err := shoppingcart.HandleEvent(msg, esdbService)
		if err != nil {
			log.Fatal(err)
		}
		msg.Ack()
	})

	a.nc.Flush()

	if err := a.nc.LastError(); err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
