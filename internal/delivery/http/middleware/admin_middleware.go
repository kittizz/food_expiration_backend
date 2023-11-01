package http_middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
)

func (m *HttpMiddleware) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := request.UserFrom(c)
		if user.Role != "admin" {
			return c.NoContent(http.StatusUnauthorized)
		}
		return next(c)
	}
}
