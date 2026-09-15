package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// DSN собирает строку подключения в формате, который понимает pgx.
func (p Configs) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.PostgresUser, p.PostgresPassword, p.PostgresHost, p.PostgresPort, p.PostgresDBName, p.PostgresSSLMode,
	)
}

type Configs struct {
	AdminToken       string `envconfig:"ADMIN_TOKEN"`
	Port             string `envconfig:"HTTP_PORT" default:"8080"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUser     string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDBName   string `envconfig:"POSTGRES_DBNAME" default:"documents"`
	PostgresSSLMode  string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
}

var Config Configs

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	Config = Configs{}
	if err := envconfig.Process("DOCUMENT", &Config); err != nil {
		fmt.Printf("ошибка конфигурации: %v", err)
		os.Exit(1)
	}
}
