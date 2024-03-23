package internal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type ShoppingCartRoutes struct {
	Service *ShoppingCartService
}

func NewShoppingCartRoutes(db *gorm.DB) *ShoppingCartRoutes {
	return &ShoppingCartRoutes{
		Service: NewShoppingCartService(db),
	}
}

func (s *ShoppingCart) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func (s *ShoppingCartRoutes) List(w http.ResponseWriter, r *http.Request) {
	shoppingCarts, _ := s.Service.List()
	list := []render.Renderer{}
	for _, v := range shoppingCarts {
		list = append(list, v)
	}
	render.RenderList(w, r, list)
}

func (s *ShoppingCartRoutes) FindOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in FindOne")
	w.WriteHeader(200)
}

type ShoppingCartRequest struct {
	*ShoppingCart
}

func (s *ShoppingCartRequest) Bind(r *http.Request) error {
	return nil
}

func (s *ShoppingCartRoutes) Create(w http.ResponseWriter, r *http.Request) {
	data := &ShoppingCartRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	s.Service.SaveOne(data.ShoppingCart.Name)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, data)
}

func (s *ShoppingCartRoutes) ShoppingCartCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var shoppingCart ShoppingCart
		var err error

		if shoppingCartID := chi.URLParam(r, "shoppingCartID"); shoppingCartID != "" {
			shoppingCartID, _ := strconv.Atoi(shoppingCartID)
			shoppingCart = s.Service.FindOne(uint(shoppingCartID))
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}

		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "shoppingCart", shoppingCart)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *ShoppingCartRoutes) Update(w http.ResponseWriter, r *http.Request) {
	shoppingCart := r.Context().Value("shoppingCart").(ShoppingCart)

	data := &ShoppingCartRequest{ShoppingCart: &shoppingCart}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err := s.Service.UpdateOne(shoppingCart.ID, &shoppingCart)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, &shoppingCart)
}

func (s *ShoppingCartRoutes) Delete(w http.ResponseWriter, r *http.Request) {
	var err error

	shoppingCart := r.Context().Value("shoppingCart").(ShoppingCart)

	err = s.Service.DeleteOne(shoppingCart.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, &shoppingCart)
}
