package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	shoppingCart "github.com/IamTossan/ecommerce_example/internal"
)

type Server struct {
	Router *chi.Mux
	DB     *gorm.DB
}

func NewDBConnection() *gorm.DB {
	dsn := "host=localhost user=user password=password dbname=ecommerce port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error creating a DB connection")
	}

	db.AutoMigrate(&shoppingCart.ShoppingCart{})

	return db
}

func NewServer() *Server {
	server := &Server{
		Router: chi.NewRouter(),
		DB:     NewDBConnection(),
	}
	return server
}

func (s *Server) MountMiddlewares() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
}

func (s *Server) MountHandlers() {
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	s.Router.Route("/shopping-carts", func(r chi.Router) {
		shoppingCartRoutes := shoppingCart.NewShoppingCartRoutes(s.DB)

		r.Get("/", shoppingCartRoutes.List)
		r.Post("/", shoppingCartRoutes.Create)

		r.Route("/{shoppingCartID}", func(r chi.Router) {
			r.Use(shoppingCartRoutes.ShoppingCartCtx)

			r.Get("/", shoppingCartRoutes.FindOne)
			r.Put("/", shoppingCartRoutes.Update)
			r.Delete("/", shoppingCartRoutes.Delete)
		})
	})

}

func main() {
	s := NewServer()
	s.MountMiddlewares()
	s.MountHandlers()

	http.ListenAndServe(":8000", s.Router)
}
