package database

import (
	_log "log"
	"os"
	"time"

	gmysql "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseMySQL struct {
	*gorm.DB
}

func NewMySQL() (*DatabaseMySQL, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}

	dsn := gmysql.NewConfig()
	dsn.User = viper.GetString("MYSQL_USERNAME")
	dsn.Passwd = viper.GetString("MYSQL_PASSWORD")
	dsn.Net = "tcp"
	dsn.Addr = viper.GetString("MYSQL_HOST") + ":" + viper.GetString("MYSQL_PORT")
	dsn.DBName = viper.GetString("MYSQL_DATABASE")
	dsn.ParseTime = true
	dsn.Loc = loc
	dsn.Collation = "utf8mb4_unicode_ci"

	instance, err := gorm.Open(
		mysql.New(mysql.Config{DSNConfig: dsn}),
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

	return &DatabaseMySQL{instance.Debug()}, nil
}
