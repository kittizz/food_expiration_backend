package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type ImageHandler struct {
	middleware   *http_middleware.HttpMiddleware
	userUsecase  domain.UserUsecase
	imageUsecase domain.ImageUsecase
}

func NewImageHandler(
	e *server.EchoServer,
	middleware *http_middleware.HttpMiddleware,
	userUsecase domain.UserUsecase,
	imageUsecase domain.ImageUsecase,
) *ImageHandler {
	handler := &ImageHandler{
		imageUsecase: imageUsecase,
		userUsecase:  userUsecase,
		middleware:   middleware,
	}

	authed := e.Group("/image", handler.middleware.AuthMiddleware)
	{
		authed.POST("/upload", handler.UploadImage)

	}
	unAuth := e.Group("/image")
	{
		unAuth.GET("/banner", handler.Getbanner)
	}
	return handler
}

func (h *ImageHandler) UploadImage(c echo.Context) error {
	user := request.UserFrom(c)

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	err = h.imageUsecase.UploadImage(c.Request().Context(), file, user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *ImageHandler) Getbanner(c echo.Context) error {
	// TODO: Implement Getbanner from database
	return c.JSON(200, echo.Map{
		"banner": "/images/banner-onlygf.png",
	})
}
