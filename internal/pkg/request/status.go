package request

import (
	"net/http"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrTokenExpired, domain.ErrInvalidDeviceId:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
