package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestMain(m *testing.M) {
	compose, err := tc.NewDockerCompose("../docker-compose.yml")
	if err != nil {
		fmt.Errorf("error in setting up docker-compose")
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = compose.Up(ctx, tc.Wait(true))
	if err != nil {
		fmt.Errorf("error in docker-compose up")
	}

	exitCode := m.Run()

	err = compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
	if err != nil {
		fmt.Errorf("error in docker-compose down")
	}
	cancel()

	os.Exit(exitCode)
}

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func afterEach(server *Server) {
	server.DB.Exec("DELETE FROM shopping_carts")
}

func TestHelloWorld(t *testing.T) {
	s := NewServer()
	s.MountHandlers()
	defer afterEach(s)

	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req, s)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "Hello World!", response.Body.String())
}

func TestShoppingCartList(t *testing.T) {
	s := NewServer()
	s.MountHandlers()
	defer afterEach(s)
	ja := jsonassert.New(t)

	req, _ := http.NewRequest("GET", "/shopping-carts", nil)
	response := executeRequest(req, s)

	require.Equal(t, http.StatusOK, response.Code)
	ja.Assertf(response.Body.String(), `[]`)
}

func TestShoppingCartCreate(t *testing.T) {
	s := NewServer()
	s.MountHandlers()
	defer afterEach(s)
	ja := jsonassert.New(t)

	req, _ := http.NewRequest("GET", "/shopping-carts", nil)
	response := executeRequest(req, s)
	ja.Assertf(response.Body.String(), `[]`)

	body := strings.NewReader(`{"name": "my-shopping-cart"}`)
	req, _ = http.NewRequest("POST", "/shopping-carts", body)
	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req, s)

	require.Equal(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("GET", "/shopping-carts", nil)
	response = executeRequest(req, s)
	ja.Assertf(
		response.Body.String(),
		`[{
			"ID": "<<PRESENCE>>",
			"name": "my-shopping-cart",
			"CreatedAt": "<<PRESENCE>>",
			"UpdatedAt": "<<PRESENCE>>",
			"DeletedAt": null
		}]`)
}

func TestShoppingCartDelete(t *testing.T) {
	s := NewServer()
	s.MountHandlers()
	defer afterEach(s)
	ja := jsonassert.New(t)

	body := strings.NewReader(`{"name": "my-shopping-cart"}`)
	req, _ := http.NewRequest("POST", "/shopping-carts", body)
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req, s)
	require.Equal(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("DELETE", "/shopping-carts/1", nil)
	response = executeRequest(req, s)
	ja.Assertf(
		response.Body.String(),
		`{
			"ID": "<<PRESENCE>>",
			"name": "my-shopping-cart",
			"CreatedAt": "<<PRESENCE>>",
			"UpdatedAt": "<<PRESENCE>>",
			"DeletedAt": null
		}`)

	req, _ = http.NewRequest("GET", "/shopping-carts", nil)
	response = executeRequest(req, s)
	ja.Assertf(response.Body.String(), `[]`)
}

func TestShoppingCartUpdate(t *testing.T) {
	s := NewServer()
	s.MountHandlers()
	defer afterEach(s)
	ja := jsonassert.New(t)

	body := strings.NewReader(`{"ID": 1, "name": "my-shopping-cart"}`)
	req, _ := http.NewRequest("POST", "/shopping-carts", body)
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req, s)
	require.Equal(t, http.StatusCreated, response.Code)

	body = strings.NewReader(`{"name": "my-shopping-cart-updated"}`)
	req, _ = http.NewRequest("PUT", "/shopping-carts/1", body)
	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req, s)
	ja.Assertf(
		response.Body.String(),
		`{
			"ID": "<<PRESENCE>>",
			"name": "my-shopping-cart-updated",
			"CreatedAt": "<<PRESENCE>>",
			"UpdatedAt": "<<PRESENCE>>",
			"DeletedAt": null
		}`)
}
