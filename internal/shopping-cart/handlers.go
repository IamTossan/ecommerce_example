package shoppingcart

import (
	"fmt"
	"net/http"
)

type ShoppingCart struct {
	Nats *NatsService
}

func (s *ShoppingCart) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in create")
	command := NewOpenShoppingCartCommand()
	s.Nats.SendCommand(r.Context(), command)
	w.WriteHeader(201)
}
func (s *ShoppingCart) Cancel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in cancel")
}
func (s *ShoppingCart) Confirm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in confirm")
}
