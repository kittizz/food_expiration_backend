package http

import (
	"time"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type LocationItemHandler struct {
	middleware *http_middleware.HttpMiddleware

	locationItemUsecase domain.LocationItemUsecase
}

func NewLocationItemHandler(e *server.EchoServer, middleware *http_middleware.HttpMiddleware, locationItemUsecase domain.LocationItemUsecase) *LocationItemHandler {
	h := &LocationItemHandler{
		middleware:          middleware,
		locationItemUsecase: locationItemUsecase,
	}

	group := e.Group("/item/:ID", h.middleware.AuthMiddleware)
	{
		group.GET("", nil)
	}
	return h
}

type createLocationItemRequest struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StorageDate time.Time `json:"storageDate"`
	ExpireDate  time.Time `json:"expireDate"`
	ForewarnDay int       `json:"forewarnDay"`
	IsArchived  bool      `json:"isArchived"`
	Category    string    `json:"category"`
	Barcode     bool      `json:"barcode"`

	ImageId    int `json:"imageId"`
	LocationId int `json:"locationId"`
}
