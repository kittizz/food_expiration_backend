package http

import (
	"net/http"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryUsercase domain.CategoryUsecase
}

func NewCategoryHandler(e *server.EchoServer, categoryUsercase domain.CategoryUsecase) *CategoryHandler {
	handler := &CategoryHandler{
		categoryUsercase: categoryUsercase,
	}

	// TODO: Add middleware admin
	unAuthed := e.Group("/category")
	{
		unAuthed.GET("", handler.ListCategory)
	}
	return handler
}

func (h *CategoryHandler) ListCategory(c echo.Context) error {
	categories, err := h.categoryUsercase.List(c.Request().Context())
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, categories)
}
