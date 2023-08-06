package http_middleware

import (
	"github.com/labstack/echo/v4/middleware"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type HttpMiddleware struct {
	userUsecase domain.UserUsecase
}

func NewHttpMiddleware(e *server.EchoServer, userUsecase domain.UserUsecase) *HttpMiddleware {
	m := &HttpMiddleware{
		userUsecase: userUsecase,
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: false,
	}))
	// e.Use(m.CORS)
	return m
}
