package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TeslaClientID     string
	TeslaClientSecret string
	TeslaRedirectURI  string
	TeslaAuthURL      string
	TeslaTokenURL     string
	TeslaAPIURL       string
	DB                struct {
		Host     string
		Port     string
		User     string
		Password string
		DbName   string
	}
	Server struct {
		Address string
	}
}

func LoadConfig() (*Config, error) {
	// 加载.env文件 如果存在
	godotenv.Load()

	cfg := &Config{}
	cfg.TeslaClientID = os.Getenv("TESLA_CLIENT_ID")
	cfg.TeslaClientSecret = os.Getenv("TESLA_CLIENT_SECRET")
	cfg.TeslaRedirectURI = os.Getenv("TESLA_REDIRECT_URI")
	cfg.TeslaAuthURL = os.Getenv("TESLA_AUTH_URL")
	cfg.TeslaTokenURL = os.Getenv("TESLA_TOKEN_URL")
	cfg.TeslaAPIURL = os.Getenv("TESLA_API_URL")

	if cfg.Server.Address == "" {
		cfg.Server.Address = ":8080"
	}
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	if cfg.DB.Port == "" {
		cfg.DB.Port = "5432"
	}
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.DbName = os.Getenv("DB_NAME")
	return cfg, nil
}
