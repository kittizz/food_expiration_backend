package http_middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
)

func (m *HttpMiddleware) UserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authorization := c.Request().Header.Get("Authorization")

		var user *domain.User
		if c.Request().Header.Get("test") == "true" || true {
			_user, err := m.userUsecase.GetAuthUserByUid(c.Request().Context(), authorization)
			if err != nil {
				return c.JSON(request.StatusCode(domain.ErrTokenExpired), request.ResponseError{Message: domain.ErrTokenExpired.Error()})
			}
			user = _user
		} else {
			_user, err := m.userUsecase.VerifyIDToken(c.Request().Context(), authorization)
			if err != nil {
				return c.JSON(request.StatusCode(domain.ErrTokenExpired), request.ResponseError{Message: domain.ErrTokenExpired.Error()})
			}
			user = _user
		}

		user, err := m.userUsecase.Sync(c.Request().Context(), *user)
		if err != nil {
			return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
		}
		return next(request.WithUser(c, user))
	}
}
