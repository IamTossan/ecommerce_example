package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"testing"

	"github.com/IamTossan/ecommerce_example/internal/application"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestMain(m *testing.M) {
	compose, err := tc.NewDockerCompose("../../docker-compose.yml")
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

func executeRequest(req *http.Request, app *application.App) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func TestBaseRoute(t *testing.T) {
	app := application.New()
	app.LoadRoutes()

	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req, app)

	require.Equal(t, http.StatusOK, response.Code)
}

func TestGetShoppingCart(t *testing.T) {
	app := application.New()
	app.LoadRoutes()

	req, _ := http.NewRequest("GET", "/shopping-carts", nil)
	response := executeRequest(req, app)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `[]`, response.Body.String())
}

func TestCreateShoppingCart(t *testing.T) {
	app := application.New()
	app.LoadRoutes()

	req, _ := http.NewRequest("GET", "/shopping-carts", nil)
	response := executeRequest(req, app)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `[]`, response.Body.String())

	req, _ = http.NewRequest("POST", "/shopping-carts", nil)
	response = executeRequest(req, app)

	require.Equal(t, http.StatusCreated, response.Code)

	req, _ = http.NewRequest("GET", "/shopping-carts", nil)
	response = executeRequest(req, app)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `[{}]`, response.Body.String())
}
