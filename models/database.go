package models

import (
	"fmt"
	"go-rest-api/config"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBName,
		config.AppConfig.DBPassword,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	if err := DB.Use(otelgorm.NewPlugin()); err != nil {
		panic("Failed to initialize OpenTelemetry GORM plugin!")
	}
	DB.AutoMigrate(&Product{})
}
