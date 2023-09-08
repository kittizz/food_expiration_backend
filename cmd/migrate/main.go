package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/kittizz/food_expiration_backend/internal/domain"
	"github.com/kittizz/food_expiration_backend/internal/pkg/database"
)

func main() {
	godotenv.Load(".env")

	viper.AutomaticEnv()

	instance, err := database.NewMySQL()
	if err != nil {
		panic(err)
	}
	if err := instance.AutoMigrate(
		&domain.User{},
		&domain.Location{},
		&domain.LocationItem{},
		&domain.ThumbnailCategory{},
		&domain.Thumbnail{},
		&domain.Blog{},
		&domain.Image{},
	); err != nil {
		panic(err)

	}
	log.Info().Msg("Migrate database")
}
