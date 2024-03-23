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

	router.Route("/shopping-carts", a.LoadShoppingCartRoutes)

	a.Router = router
}

func (a *App) LoadShoppingCartRoutes(router chi.Router) {
	shoppingCartHandler := &shoppingcart.ShoppingCart{
		Nats: &shoppingcart.NatsService{
			Nc: a.nc,
		},
		Esdb: &shoppingcart.EsdbService{
			Esdb: a.es,
		},
	}

	router.Get("/", shoppingCartHandler.List)
	router.Post("/", shoppingCartHandler.Create)
	router.Post("/{id}/cancel", shoppingCartHandler.Cancel)
	router.Post("/{id}/confirm", shoppingCartHandler.Confirm)
}
