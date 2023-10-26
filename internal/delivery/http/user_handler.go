package http

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/messaging"
	"github.com/labstack/echo/v4"

	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/firebase"
	"github.com/kittizz/food_expiration_backend/internal/pkg/request"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

type UserHandler struct {
	middleware   *http_middleware.HttpMiddleware
	userUsecase  domain.UserUsecase
	imageUsecase domain.ImageUsecase
	firebase     *firebase.Firebase
}

func NewUserHandler(
	e *server.EchoServer,
	middleware *http_middleware.HttpMiddleware,
	userUsecase domain.UserUsecase,
	imageUsecase domain.ImageUsecase,
	firebase *firebase.Firebase,

) *UserHandler {
	handler := &UserHandler{
		userUsecase:  userUsecase,
		middleware:   middleware,
		imageUsecase: imageUsecase,
		firebase:     firebase,
	}
	unAuth := e.Group("/user")
	{
		unAuth.POST("/register-device", handler.RegisteDevicer)

	}

	authed := e.Group("/user", handler.middleware.AuthMiddleware)
	{
		authed.GET("", handler.GetUser)
		authed.POST("/change-profilepicture", handler.ChangeProfilePicture)
		authed.POST("/change-nickname", handler.ChangeNickname)
		authed.POST("/update-fcm", handler.UpdateFcm)
		authed.POST("/update-notifications", handler.UpdateNotifications)

	}
	e.GET("/test_token", handler.TestToken)
	e.POST("/test_fcm", handler.TestFcm)

	return handler
}

type RegisteDevicerRequest struct {
	AuthToken string `json:"auth_token"`

	Nickname string `json:"nickname"`
}

func (h *UserHandler) RegisteDevicer(c echo.Context) error {

	var req RegisteDevicerRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := h.userUsecase.VerifyIDToken(c.Request().Context(), req.AuthToken)
	if err != nil {
		log.Print(err)
		return c.JSON(request.StatusCode(domain.ErrTokenExpired), request.ResponseError{Message: domain.ErrTokenExpired.Error()})
	}

	user, err = h.userUsecase.Sync(c.Request().Context(), *user)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	deviceId, err := h.userUsecase.RegisterDevice(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	user.DeviceId = &deviceId

	return c.JSON(200, echo.Map{
		"deviceId": deviceId,
	})
}

func (h *UserHandler) TestToken(c echo.Context) error {
	token, err := h.userUsecase.GenerateIDToken(c.Request().Context(), "QCBlS58JmMSZDQyednJYKWgeShi1")
	if err != nil {
		return err
	}
	return c.String(200, token)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	user := *request.UserFrom(c)

	if user.ProfilePictureBlurHash == nil || user.ProfilePicture == nil {
		str := "LIEpzCa#1mt7EjWB?Hof5Xoe}fR%"
		str2 := "/images/user.png"
		user.ProfilePictureBlurHash = &str
		user.ProfilePicture = &str2
	}
	return c.JSON(200, user)
}

func (h *UserHandler) ChangeProfilePicture(c echo.Context) error {
	user := request.UserFrom(c)

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	err = h.userUsecase.ChangeProfile(c.Request().Context(), file, c.FormValue("hash"), user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	if user.ProfilePicture != nil {
		err := h.imageUsecase.DeleteWithPath(c.Request().Context(), *user.ProfilePicture)
		if err != nil {
			return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
		}
	}

	return c.NoContent(http.StatusOK)
}

type changeNicknameRequest struct {
	Nickname string `json:"nickname"`
}

func (h *UserHandler) ChangeNickname(c echo.Context) error {
	var req changeNicknameRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user := request.UserFrom(c)
	err := h.userUsecase.ChangeNickname(c.Request().Context(), req.Nickname, user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

type updateFcmRequest struct {
	FcmToken   *string `json:"fcmToken"`
	DeviceType *string `json:"deviceType"`
}

func (h *UserHandler) UpdateFcm(c echo.Context) error {
	var req updateFcmRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user := request.UserFrom(c)
	err := h.userUsecase.UpdateFcm(c.Request().Context(), req.FcmToken, req.DeviceType, user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

type updateNotificationsRequest struct {
	Notifications bool `json:"notifications"`
}

func (h *UserHandler) UpdateNotifications(c echo.Context) error {
	var req updateNotificationsRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	user := request.UserFrom(c)
	err := h.userUsecase.UpdateNotification(c.Request().Context(), &req.Notifications, user.ID)
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
func (h *UserHandler) TestFcm(c echo.Context) error {
	f, err := h.firebase.FcmClient()
	if err != nil {
		return c.JSON(request.StatusCode(err), request.ResponseError{Message: err.Error()})
	}
	f.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Congratulations!!",
			Body:  "You have just implement push notification",
		},
		Token: c.FormValue("token"), // it's a single device token

	})
	return c.NoContent(http.StatusOK)
}
