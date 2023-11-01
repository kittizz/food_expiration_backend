package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type AdminHandler struct {
	adminUsecase domain.AdminUsecase
}

func NewAdminHandler(e *server.EchoServer, adminUsecase domain.AdminUsecase) *AdminHandler {
	h := &AdminHandler{
		adminUsecase: adminUsecase,
	}
	adminGroup := e.Group("/admin")
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
