package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type ItemHandler struct {
	middleware *http_middleware.HttpMiddleware

	itemUsecase domain.ItemUsecase
}

func NewItemHandler(e *server.EchoServer, middleware *http_middleware.HttpMiddleware, itemUsecase domain.ItemUsecase) *ItemHandler {
	h := &ItemHandler{
		middleware:  middleware,
		itemUsecase: itemUsecase,
	}

	group := e.Group("/item", h.middleware.AuthMiddleware)
	{
		group.GET("/:id", h.GetItem)
		group.POST("/create", h.CreateItem)
		group.GET("/location", h.GetLocationItemAll)
		group.GET("/location/:id", h.GetLocationItem)
		group.PUT("/clear", h.ClearItem)
		group.PUT("/:id", h.UpdateItem)

	}
	return h
}

func (h *ItemHandler) GetItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	item, err := h.itemUsecase.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

type createItemRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StorageDate time.Time `json:"storageDate"`
	ExpireDate  time.Time `json:"expireDate"`
	ForewarnDay int       `json:"forewarnDay"`
	IsArchived  bool      `json:"isArchived"`
	Category    string    `json:"category"`
	Barcode     string    `json:"barcode"`

	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`

	ImageId    int `json:"imageId"`
	LocationId int `json:"locationId"`
}

func (h *ItemHandler) CreateItem(c echo.Context) error {
	var req createItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	err := h.itemUsecase.Create(c.Request().Context(), domain.Item{
		Name:        &req.Name,
		Description: &req.Description,
		StorageDate: req.StorageDate,
		ExpireDate:  req.ExpireDate,
		ForewarnDay: &req.ForewarnDay,
		IsArchived:  &req.IsArchived,
		Category:    &req.Category,
		Barcode:     &req.Barcode,
		Quantity:    &req.Quantity,
		Unit:        &req.Unit,

		ImageID:    req.ImageId,
		LocationID: req.LocationId,
	})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusCreated)
}

type updateItemRequest struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StorageDate time.Time `json:"storageDate"`
	ExpireDate  time.Time `json:"expireDate"`
	ForewarnDay int       `json:"forewarnDay"`
	IsArchived  bool      `json:"isArchived"`
	Category    string    `json:"category"`
	Barcode     string    `json:"barcode"`

	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`

	ImageId    int `json:"imageId"`
	LocationId int `json:"locationId"`
}

func (h *ItemHandler) UpdateItem(c echo.Context) error {
	var req updateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	fmt.Println(req.StorageDate)
	err := h.itemUsecase.UpdateByID(c.Request().Context(), domain.Item{
		Name:        &req.Name,
		Description: &req.Description,
		StorageDate: req.StorageDate,
		ExpireDate:  req.ExpireDate,
		ForewarnDay: &req.ForewarnDay,
		IsArchived:  &req.IsArchived,
		Category:    &req.Category,
		Barcode:     &req.Barcode,
		Quantity:    &req.Quantity,
		Unit:        &req.Unit,

		ImageID:    req.ImageId,
		LocationID: req.LocationId,
	}, req.Id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *ItemHandler) GetLocationItem(c echo.Context) error {
	var id *int
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	if idInt <= 0 {
		id = nil
	} else {
		id = &idInt
	}

	isArchivedBool, err := strconv.ParseBool(c.QueryParam("isArchived"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	item, err := h.itemUsecase.List(c.Request().Context(), id, isArchivedBool, true)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}
func (h *ItemHandler) GetLocationItemAll(c echo.Context) error {
	isArchivedBool, err := strconv.ParseBool(c.QueryParam("isArchived"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	items, err := h.itemUsecase.List(c.Request().Context(), nil, isArchivedBool, true)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

type clearItemRequest struct {
	Id      []int `json:"id"`
	Archive bool  `json:"archive"`
}

func (h *ItemHandler) ClearItem(c echo.Context) error {

	var req clearItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err := h.itemUsecase.Archive(c.Request().Context(), req.Archive, req.Id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
