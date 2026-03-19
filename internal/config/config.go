package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	Port             string
}

func (c *Config) GetMissingFields() []string {
	var missing []string

	check := []struct {
		value string
		name  string
	}{
		{c.DatabaseHost, "DATABASE_HOST"},
		{c.DatabasePort, "DATABASE_PORT"},
		{c.DatabaseName, "DATABASE_NAME"},
		{c.DatabaseUser, "DATABASE_USER"},
		{c.DatabasePassword, "DATABASE_PASSWORD"},
		{c.Port, "PORT"},
	}

	// go al recorrer un array dinamico (slice) devuelve indice y valor, como no usamos indice no podemos hacer una nomenclatura para no usarla, por eso, en lugar de algo como index o i usamos _ que es como un blank identifier
	for _, field := range check {
		if field.value == "" {
			missing = append(missing, field.name)
		}
	}
	return missing
}

func NewDatabase(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName,
	)

	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db, nil
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Println("Info: .env file not found, using enviroment variables")
	}

	config := &Config{
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		Port:             os.Getenv("PORT"),
	}

	missing := config.GetMissingFields()
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}

	return config, nil
}
