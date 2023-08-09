package http

import (
	"log"

	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type UserHandler struct {
	middleware  *http_middleware.HttpMiddleware
	userUsecase domain.UserUsecase
}

func NewUserHandler(
	e *server.EchoServer,
	middleware *http_middleware.HttpMiddleware,
	userUsecase domain.UserUsecase,

) *UserHandler {
	handler := &UserHandler{
		userUsecase: userUsecase,
		middleware:  middleware,
	}
	unAuth := e.Group("/user")
	{
		unAuth.POST("/register-device", handler.RegisteDevicer)

	}

	authed := e.Group("/user", handler.middleware.AuthMiddleware)
	{
		authed.GET("", handler.GetUser)

	}
	e.GET("/test_token", handler.TestToken)
	return handler
}

type RegisteDevicerRequest struct {
	AuthToken string `json:"auth_token"`
}

func (h *UserHandler) RegisteDevicer(c echo.Context) error {

	var req RegisteDevicerRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := h.userUsecase.VerifyIDToken(c.Request().Context(), req.AuthToken)
	if err != nil {
		log.Print(err)
		return c.JSON(request.StatusCode(domain.ErrTokenExpired), request.ResponseError{Message: domain.ErrTokenExpired.Error()})
	}

	user, err = h.userUsecase.Sync(c.Request().Context(), *user)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	deviceId, err := h.userUsecase.RegisterDevice(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	user.DeviceId = &deviceId

	return c.JSON(200, echo.Map{
		"deviceId": deviceId,
	})
}

func (h *UserHandler) TestToken(c echo.Context) error {
	token, err := h.userUsecase.GenerateIDToken(c.Request().Context(), "QCBlS58JmMSZDQyednJYKWgeShi1")
	if err != nil {
		return err
	}
	return c.String(200, token)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	user := request.UserFrom(c)
	return c.JSON(200, user)
}
