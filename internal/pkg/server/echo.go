package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type EchoServer struct {
	*echo.Echo
}

func NewEchoServer(lc fx.Lifecycle) *EchoServer {
	e := echo.New()

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
