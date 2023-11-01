package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type AdminHandler struct {
	adminUsecase    domain.AdminUsecase
	categoryUsecase domain.CategoryUsecase
	middleware      *http_middleware.HttpMiddleware
}

func NewAdminHandler(e *server.EchoServer, adminUsecase domain.AdminUsecase, middleware *http_middleware.HttpMiddleware, categoryUsecase domain.CategoryUsecase) *AdminHandler {
	h := &AdminHandler{
		adminUsecase:    adminUsecase,
		categoryUsecase: categoryUsecase,
		middleware:      middleware,
	}
	// TODO: remove this once
	// adminGroup := e.Group("/admin", h.middleware.AuthMiddleware, h.middleware.AdminMiddleware)
	adminGroup := e.Group("/admin")
	{
		adminGroup.GET("/dashboard", h.GetDashboard)
		adminGroup.GET("/category", h.GetCategory)
		adminGroup.POST("/set-category", h.SetCategory)
	}
	return h
}

func (h *AdminHandler) GetDashboard(c echo.Context) error {
	dashboard, err := h.adminUsecase.Dashboard(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dashboard)
}
func (h *AdminHandler) GetCategory(c echo.Context) error {
	categories, err := h.categoryUsecase.List(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"categories": strings.Join(categories, "\n")})
}

type setCategoryRequest struct {
	Categories string `json:"categories"`
}

func (h *AdminHandler) SetCategory(c echo.Context) error {
	var req setCategoryRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	err := h.categoryUsecase.Set(c.Request().Context(), req.Categories)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
