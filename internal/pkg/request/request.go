package request

import (
	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

func WithUser(c echo.Context, user *domain.User) echo.Context {
	c.Set("$User", user)
	return c
}

func UserFrom(c echo.Context) *domain.User {
	return c.Get("$User").(*domain.User)
}
