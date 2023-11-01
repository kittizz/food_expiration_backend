package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type AdminHandler struct {
	adminUsecase domain.AdminUsecase
	middleware   *http_middleware.HttpMiddleware
}

func NewAdminHandler(e *server.EchoServer, adminUsecase domain.AdminUsecase, middleware *http_middleware.HttpMiddleware) *AdminHandler {
	h := &AdminHandler{
		adminUsecase: adminUsecase,
		middleware:   middleware,
	}
	adminGroup := e.Group("/admin", h.middleware.AuthMiddleware, h.middleware.AdminMiddleware)
	{
		adminGroup.GET("/dashboard", h.GetDashboard)
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
