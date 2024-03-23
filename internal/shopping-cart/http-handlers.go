package shoppingcart

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ShoppingCart struct {
	Nats *NatsService
	Esdb *EsdbService
}

func (s *ShoppingCart) List(w http.ResponseWriter, r *http.Request) {
	render.RenderList(w, r, []render.Renderer{})
}

func (s *ShoppingCart) Create(w http.ResponseWriter, r *http.Request) {
	command := NewOpenShoppingCartCommand()
	s.Nats.SendCommand(r.Context(), command)
	w.WriteHeader(201)
}

func (s *ShoppingCart) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id was not UUID type"))
		return
	}
	command := NewCancelShoppingCartCommand(uuid)
	s.Nats.SendCommand(r.Context(), command)
	w.WriteHeader(200)
}

func (s *ShoppingCart) Confirm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in confirm")
}
