package http_middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
)

func (m *HttpMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// deviceId := c.Request().URL.Query().Get("deviceId")
		deviceId := c.Request().Header.Get("x-device-id")
		user, err := m.userUsecase.GetUserByDeviceId(c.Request().Context(), deviceId)
		if err != nil {
			return c.JSON(request.StatusCode(domain.ErrInvalidDeviceId), request.ResponseError{Message: domain.ErrInvalidDeviceId.Error()})
		}
		return next(request.WithUser(c, user))
	}
}
