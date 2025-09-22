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
	Server            struct {
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
	return cfg, nil
}
