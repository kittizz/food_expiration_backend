package http

import (
	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type TestHandler struct {
}

func NewTestHandler(e *server.EchoServer) *TestHandler {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World")
	})
	return &TestHandler{}
}
