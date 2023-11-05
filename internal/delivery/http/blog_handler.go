package http

import (
	"fmt"
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

		//TODO: add middleware
		authed := blogGroup.Group("")
		{
			authed.POST("/rename", h.Rename)
			authed.PUT("/update", h.Update)
			authed.DELETE("", h.Delete)
		}
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

type renameBlogRequest struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func (h *BlogHandler) Rename(c echo.Context) error {
	var req renameBlogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err := h.blogUsecase.Rename(c.Request().Context(), req.Name, req.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

type updateBlogRequest struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	ImageId int    `json:"imageId"`
}

func (h *BlogHandler) Update(c echo.Context) error {
	var req updateBlogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	id, err := h.blogUsecase.UpdateOrCreate(c.Request().Context(), domain.Blog{
		Title:   req.Title,
		Content: req.Content,
		ImageID: req.ImageId,
	}, req.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.String(http.StatusOK, fmt.Sprintf("%v", id))
}
func (h *BlogHandler) Delete(c echo.Context) error {
	id := c.QueryParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	err = h.blogUsecase.Delete(c.Request().Context(), idInt)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusOK)
}
