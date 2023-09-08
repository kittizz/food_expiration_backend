package backend

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
	"github.com/kittizz/food_expiration_backend/internal/pkg/auth"
	"github.com/kittizz/food_expiration_backend/internal/pkg/bucket"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
	"github.com/kittizz/food_expiration_backend/internal/pkg/server"
	"github.com/kittizz/food_expiration_backend/internal/repository"
	"github.com/kittizz/food_expiration_backend/internal/usecase"
)

func init() {

	godotenv.Load(".env")

	viper.AutomaticEnv()

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)

	if viper.GetString("LEVEL") == "DEBUG" {
		log.Print("DEBUG ON")
		log.Level(zerolog.DebugLevel)
	}

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Panic().Msg("cannot load location")
	}
	time.Local = loc
}

func Run(firebaseCredentials []byte) {
	app := fx.New(
		fx.Provide(
			database.NewMySQL,
			repository.NewUserRepository,
			repository.NewLocationRepository,
			repository.NewBlogRepository,
			repository.NewImageRepository,

			auth.NewFirebase(firebaseCredentials),
			bucket.NewBucket,
			server.NewEchoServer,

			usecase.NewUserUsecase,
			usecase.NewLocationUsecase,
			usecase.NewBlogUsecase,
			usecase.NewImageUsecase,

			http_middleware.NewHttpMiddleware,

			http.NewUserHandler,
			http.NewBlogHandler,
			http.NewLocationHandler,
			http.NewImageHandler,
		),

		fx.Invoke(func(*http.UserHandler, *http.BlogHandler, *http.LocationHandler, *http.ImageHandler) {
		}),
	)

	app.Run()

}
