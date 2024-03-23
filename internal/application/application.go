package application

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/nats-io/nats.go"
)

type App struct {
	env    *Environment
	Router http.Handler
	nc     *nats.Conn
	es     *esdb.Client
}

func New() *App {
	env := NewEnvironment()
	app := &App{
		env: env,
		nc:  GetNatsConnection(env.NatsUrl),
		es:  GetEventStoreConnection(env.EsdbUrl),
	}
	return app
}

func (a *App) StartServer(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.env.Port,
		Handler: a.Router,
	}

	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		defer a.nc.Drain()
		log.Println("draining nats and shutting down")

		return server.Shutdown(context.Background())
	}
}
