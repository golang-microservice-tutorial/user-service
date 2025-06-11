package config

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	encodedPassword := url.QueryEscape(Config.Database.Password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		Config.Database.Username,
		encodedPassword,
		Config.Database.Host,
		Config.Database.Port,
		Config.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(Config.Database.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(Config.Database.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.Database.MaxLifeTimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(Config.Database.MaxidleTime) * time.Second)

	return db, nil
}
