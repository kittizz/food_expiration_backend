package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type BlogHandler struct {
	blogUsecase domain.BlogUsecase
}

func NewBlogHandler(e *server.EchoServer, blogUsecase domain.BlogUsecase) *BlogHandler {
	h := &BlogHandler{
		blogUsecase: blogUsecase,
	}
	blogGroup := e.Group("/blog")
	{
		blogGroup.GET("/recommend", h.GetRecommend)
		blogGroup.GET("/all", h.GetList)
		blogGroup.GET("/query", h.GetByID)
	}
	return h
}

func (h *BlogHandler) GetRecommend(c echo.Context) error {
	blogs, err := h.blogUsecase.List(c.Request().Context(), true, 3)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, blogs)
}

func (h *BlogHandler) GetList(c echo.Context) error {
	blogs, err := h.blogUsecase.List(c.Request().Context(), false, 0)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, blogs)
}
func (h *BlogHandler) GetByID(c echo.Context) error {
	id := c.QueryParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	blog, err := h.blogUsecase.GetByID(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, blog)
}
