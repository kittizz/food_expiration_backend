package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ziflex/lecho/v3"
	"go.uber.org/fx"
)

type EchoServer struct {
	*echo.Echo
}

func NewEchoServer(lc fx.Lifecycle) *EchoServer {
	e := echo.New()
	e.HideBanner = true

	elog := lecho.From(log.Logger)
	e.Logger = elog
	// e.Use(middleware.RequestID())
	e.Use(lecho.Middleware(lecho.Config{
		Logger: elog,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info().Msgf("Starting echo server on port %s", viper.GetString("HTTP_PORT"))
				if err := e.Start(":" + viper.GetString("HTTP_PORT")); err != nil && err != http.ErrServerClosed {
					log.Info().Msgf("echo.Start error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Stopping echo server")
			return e.Shutdown(ctx)
		},
	})

	return &EchoServer{e}
}
