package application

import (
	"net/http"

	shoppingcart "github.com/IamTossan/ecommerce_example/internal/shopping-cart"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) LoadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	router.Route("/shopping-cart", a.LoadShoppingCartRoutes)

	a.router = router
}

func (a *App) LoadShoppingCartRoutes(router chi.Router) {
	shoppingCartHandler := &shoppingcart.ShoppingCart{
		Nats: &shoppingcart.NatsService{
			Nc: a.nc,
		},
	}

	router.Post("/", shoppingCartHandler.Create)
	router.Post("/{id}/cancel", shoppingCartHandler.Cancel)
	router.Post("/{id}/confirm", shoppingCartHandler.Confirm)
}
