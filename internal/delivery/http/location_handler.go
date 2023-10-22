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

		group.GET("", h.GetLocation)
		group.GET("/list", h.GetLocationList)
		group.POST("", h.CreateLocation)
		group.DELETE("", h.DeleteLocation)
		group.PUT("", h.UpdateLocation)
	}
	return h
}

func (h *LocationHandler) GetLocation(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}

	items := c.QueryParam("items") == "true"

	location, err := h.locationUsecase.Get(c.Request().Context(), id, items)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, location)
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
	ImageId     int    `json:"imageId"`
}

func (h *LocationHandler) CreateLocation(c echo.Context) error {
	user := request.UserFrom(c)
	var req createLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	location := domain.Location{
		Name:        &req.Name,
		Description: &req.Description,
		ImageID:     req.ImageId,
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

type updateLocationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageId     int    `json:"imageId"`
	LocationId  int    `json:"locationId"`
}

func (h *LocationHandler) UpdateLocation(c echo.Context) error {
	var req updateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	location := domain.Location{
		Name:        &req.Name,
		Description: &req.Description,
		ImageID:     req.ImageId,
	}
	if err := h.locationUsecase.UpdateByID(c.Request().Context(), location, req.LocationId); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)

}
