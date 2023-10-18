package http

import (
	"net/http"
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
		group.POST("/create", h.CreateItem)

	}
	return h
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

	ImageId    int `json:"imageId"`
	LocationId int `json:"locationId"`
}

func (h *ItemHandler) CreateItem(c echo.Context) error {
	var req createItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	err := h.itemUsecase.Create(c.Request().Context(), domain.Item{
		Name:        req.Name,
		Description: req.Description,
		StorageDate: req.StorageDate,
		ExpireDate:  req.ExpireDate,
		ForewarnDay: req.ForewarnDay,
		IsArchived:  req.IsArchived,
		Category:    req.Category,
		Barcode:     req.Barcode,

		ImageID:    req.ImageId,
		LocationID: req.LocationId,
	})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusCreated)
}
