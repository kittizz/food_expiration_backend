package http_middleware

import (
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type HttpMiddleware struct {
}

func NewHttpMiddleware(e *server.EchoServer) *HttpMiddleware {
	m := &HttpMiddleware{}

	e.Use(m.CORS)
	return m
}
