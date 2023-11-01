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

type ThumbnailHandler struct {
	middleware *http_middleware.HttpMiddleware

	thumbnailCategoryUsecase domain.ThumbnailCategoryUsecase
	thumbnailUsecase         domain.ThumbnailUsecase
}

func NewThumbnailHandler(e *server.EchoServer, middleware *http_middleware.HttpMiddleware, thumbnailCategoryUsecase domain.ThumbnailCategoryUsecase, thumbnailUsecase domain.ThumbnailUsecase) *ThumbnailHandler {
	h := &ThumbnailHandler{
		middleware:               middleware,
		thumbnailCategoryUsecase: thumbnailCategoryUsecase,
		thumbnailUsecase:         thumbnailUsecase,
	}

	//TODO: add middleware

	group := e.Group("/thumbnail")
	{
		group.GET("", h.GetThumbnailByCatrgoryId)
		group.POST("/create-thumbnail", h.CreateThumbnail)
		group.DELETE("", h.DeleteThumbnail)

		group.GET("/category", h.GetCategory)
		group.POST("/create-category", h.CreateCategory)
		group.DELETE("/category", h.DeleteCategory)
		group.PUT("/update-category-image", h.UpdateCategoryImage)

		group.POST("/rename", h.Rename)
	}
	return h
}

type createThumbnailRequest struct {
	Name       string `json:"name"`
	ImageID    int    `json:"imageId"`
	CategoryID int    `json:"categoryId"`
}

func (h *ThumbnailHandler) CreateThumbnail(c echo.Context) error {
	var req createThumbnailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	thumbnail := domain.Thumbnail{
		Name:                req.Name,
		ImageID:             req.ImageID,
		ThumbnailCategoryID: req.CategoryID,
	}

	if err := h.thumbnailUsecase.Create(c.Request().Context(), &thumbnail); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, nil)
}
func (h *ThumbnailHandler) GetCategory(c echo.Context) error {
	list, err := h.thumbnailCategoryUsecase.List(c.Request().Context())
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, list)

}

type createThumbnailCategoryRequest struct {
	Name    string `json:"name"`
	ImageID int    `json:"imageId"`
	Type    string `json:"type"`
}

func (h *ThumbnailHandler) CreateCategory(c echo.Context) error {
	var req createThumbnailCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	category := domain.ThumbnailCategory{
		Name:    req.Name,
		ImageID: req.ImageID,
		Type:    req.Type,
	}
	if err := h.thumbnailCategoryUsecase.Create(c.Request().Context(), &category); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, nil)
}

func (h *ThumbnailHandler) DeleteCategory(c echo.Context) error {
	idInt, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err = h.thumbnailUsecase.Delete(c.Request().Context(), domain.Thumbnail{ThumbnailCategoryID: idInt})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	err = h.thumbnailCategoryUsecase.Delete(c.Request().Context(), domain.ThumbnailCategory{
		ID: idInt,
	})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
func (h *ThumbnailHandler) GetThumbnailByCatrgoryId(c echo.Context) error {
	idInt, err := strconv.Atoi(c.QueryParam("catrgoryId"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	thumbnail, err := h.thumbnailCategoryUsecase.Get(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, thumbnail)
}

type renameRequest struct {
	Id int `json:"id"`

	Name string `json:"name"`
	Type string `json:"type"`
}

func (h *ThumbnailHandler) Rename(c echo.Context) error {
	var req renameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	if req.Type == "category" {
		err := h.thumbnailCategoryUsecase.Update(c.Request().Context(), domain.ThumbnailCategory{
			Name: req.Name,
		}, req.Id)
		if err != nil {
			return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
		}
	}
	if req.Type == "thumbnail" {
		err := h.thumbnailUsecase.Update(c.Request().Context(), domain.Thumbnail{
			Name: req.Name,
		}, req.Id)
		if err != nil {
			return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ThumbnailHandler) DeleteThumbnail(c echo.Context) error {
	idInt, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err = h.thumbnailUsecase.Delete(c.Request().Context(), domain.Thumbnail{
		ID: idInt,
	})
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

type updateCategoryImageRequest struct {
	Id      int `json:"id"`
	ImageID int `json:"imageId"`
}

func (h *ThumbnailHandler) UpdateCategoryImage(c echo.Context) error {
	var req updateCategoryImageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err := h.thumbnailCategoryUsecase.Update(c.Request().Context(), domain.ThumbnailCategory{
		ImageID: req.ImageID,
	}, req.Id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
