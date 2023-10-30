package http

import (
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
		group.POST("/location", h.GetLocationItemAll)
		group.PUT("/clear", h.ClearItem)
		group.PUT("/:id", h.UpdateItem)
		group.PATCH("/delete", h.DeleteItem)

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

	user := request.UserFrom(c)
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
		UserID:     user.ID,
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

	ImageId           int  `json:"imageId"`
	LocationId        int  `json:"locationId"`
	ResetNotification bool `json:"resetNotification"`
}

func (h *ItemHandler) UpdateItem(c echo.Context) error {
	var req updateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	user := request.UserFrom(c)

	item := domain.Item{
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
		ImageID:     req.ImageId,
		LocationID:  req.LocationId,
		UserID:      user.ID,
	}
	if req.ResetNotification {
		now, _ := time.Parse(time.DateOnly, time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
		item.LastNotificationAt = &now
		item.NotificationStatus = domain.NOTIFICATION_STATUS_PLANNED
	}
	err := h.itemUsecase.UpdateByID(c.Request().Context(), item, req.Id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

type locationItemAll struct {
	Id         int  `json:"id"`
	IsArchived bool `json:"isArchived"`
}

func (h *ItemHandler) GetLocationItemAll(c echo.Context) error {
	var req locationItemAll
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	user := request.UserFrom(c)

	items, err := h.itemUsecase.List(c.Request().Context(), user.ID, req.Id, req.IsArchived, true)
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

type deleteItemRequest struct {
	Id []int `json:"id"`
}

func (h *ItemHandler) DeleteItem(c echo.Context) error {
	var req deleteItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err := h.itemUsecase.Deletes(c.Request().Context(), req.Id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
