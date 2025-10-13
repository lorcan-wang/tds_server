package config

import (
	"os"
	"path/filepath"

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
	loadEnv()

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

// loadEnv attempts to load environment variables from a .env file.
// It looks in the current working directory first, then walks up to parent directories.
func loadEnv() {
	if err := godotenv.Load(); err == nil {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		envPath := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(envPath); statErr == nil {
			_ = godotenv.Load(envPath)
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
	}
}
