package database

import (
	"fmt"
	_log "log"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kittizz/food_expiration_backend/internal/domain"
)

type DatabaseMySQL struct {
	*gorm.DB
}

func NewMySQL() (*DatabaseMySQL, error) {
	instance, err := gorm.Open(
		mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s",
			viper.GetString("MYSQL_USERNAME"),
			viper.GetString("MYSQL_PASSWORD"),
			viper.GetString("MYSQL_HOST"),
			viper.GetInt("MYSQL_PORT"),
			viper.GetString("MYSQL_DATABASE"),
			"Asia%2FBangkok",
		)),
		&gorm.Config{
			Logger: logger.New(
				_log.New(os.Stdout, "\n", _log.LstdFlags),
				logger.Config{
					LogLevel: logger.Silent,
				},
			),
		},
	)
	if err != nil {
		return nil, err
	}

	db, err := instance.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	if viper.GetBool("MYSQL_MIGRATE") {
		if err := instance.AutoMigrate(
			&domain.User{},
			&domain.Location{},
			&domain.LocationItem{},
			&domain.ThumbnailCategory{},
			&domain.Thumbnail{},
			&domain.Blog{},
		); err != nil {
			return nil, err
		}
		log.Info().Msg("Migrate database")
		os.Exit(1)
	}
	return &DatabaseMySQL{instance.Debug()}, nil
}
