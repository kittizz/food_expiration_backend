package http

import (
	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type UserHandler struct {
	middleware      *http_middleware.HttpMiddleware
	userUsecase     domain.UserUsecase
	locationUsecase domain.LocationUsecase
}

func NewUserHandler(
	e *server.EchoServer,
	middleware *http_middleware.HttpMiddleware,
	userUsecase domain.UserUsecase,
	locationUsecase domain.LocationUsecase,
) *UserHandler {
	handler := &UserHandler{
		userUsecase:     userUsecase,
		middleware:      middleware,
		locationUsecase: locationUsecase,
	}
	group := e.Group("/user", middleware.UserAuth)
	{
		group.GET("", handler.GetUser)

	}
	e.GET("/test_token", handler.TestToken)
	return handler
}

func (h *UserHandler) GetUser(c echo.Context) error {
	user := request.UserFrom(c)
	return c.JSON(200, user)
}

func (h *UserHandler) TestToken(c echo.Context) error {
	token, err := h.userUsecase.GenerateIDToken(c.Request().Context(), "QCBlS58JmMSZDQyednJYKWgeShi1")
	if err != nil {
		return err
	}
	return c.String(200, token)
}
