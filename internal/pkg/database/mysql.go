package database

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQL() (*gorm.DB, error) {
	instance, err := gorm.Open(
		mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4_unicode_ci&parseTime=true&loc=%s",
			viper.GetString("MYSQL_USERNAME"),
			viper.GetString("MYSQL_PASSWORD"),
			viper.GetString("MYSQL_HOST"),
			viper.GetInt("MYSQL_PORT"),
			viper.GetString("MYSQL_DATABASE"),
			"Asia%2FBangkok",
		)),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\n", log.LstdFlags),
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
		// &User{},
		); err != nil {
			return nil, err
		}
	}
	return instance, nil
}
