package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type LocationHandler struct {
	middleware *http_middleware.HttpMiddleware

	locationUsecase domain.LocationUsecase
}

func NewLocationHandler(e *server.EchoServer, middleware *http_middleware.HttpMiddleware, locationUsecase domain.LocationUsecase) *LocationHandler {
	h := &LocationHandler{
		middleware:      middleware,
		locationUsecase: locationUsecase,
	}

	group := e.Group("/location", h.middleware.AuthMiddleware)
	{
		group.GET("/list", h.GetLocationList)
		group.POST("", h.CreateLocation)
		group.DELETE("", h.DeleteLocation)
	}
	return h
}

func (h *LocationHandler) GetLocationList(c echo.Context) error {
	user := request.UserFrom(c)
	locations, err := h.locationUsecase.List(c.Request().Context(), domain.Location{
		UserID: user.ID,
	})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, locations)
}

type createLocationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func (h *LocationHandler) CreateLocation(c echo.Context) error {
	user := request.UserFrom(c)
	var req createLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	location := domain.Location{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		UserID:      user.ID,
	}
	if err := h.locationUsecase.Create(c.Request().Context(), location); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, nil)
}

func (h *LocationHandler) DeleteLocation(c echo.Context) error {
	idInt, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err = h.locationUsecase.Delete(c.Request().Context(), domain.Location{
		ID: idInt,
	})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}
