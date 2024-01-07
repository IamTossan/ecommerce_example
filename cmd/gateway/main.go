package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/IamTossan/ecommerce_example/internal/application"
)

func main() {
	app := application.New()
	app.LoadRoutes()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app.StartServer(ctx)
}
