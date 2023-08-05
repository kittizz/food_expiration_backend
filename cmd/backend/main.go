package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/kittizz/food_expiration_backend/internal/delivery/http"
	http_middleware "github.com/kittizz/food_expiration_backend/internal/delivery/http/middleware"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	log.Logger = log.Output(output)

	godotenv.Load(".env")

	viper.AutomaticEnv()

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Panic().Msg("cannot load location")
	}
	time.Local = loc
}

func main() {
	app := fx.New(
		fx.Provide(
			database.NewMySQL,

			server.NewEchoServer,
			http_middleware.NewHttpMiddleware,
			http.NewTestHandler,
		),

		fx.Invoke(func(*http.TestHandler) {}),
	)

	app.Run()

}
