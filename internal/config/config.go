package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	TeslaClientID        string
	TeslaClientSecret    string
	TeslaRedirectURI     string
	TeslaAuthURL         string
	TeslaTokenURL        string
	TeslaAPIURL          string
	TeslaPartnerTokenURL string
	TeslaPartnerScope    string
	TeslaPartnerDomain   string
	TeslaCommandKeyPath  string
	DB                   struct {
		Host     string
		Port     string
		User     string
		Password string
		DbName   string
		TimeZone string
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
	if cfg.TeslaAPIURL == "" {
		cfg.TeslaAPIURL = "https://fleet-api.prd.cn.vn.cloud.tesla.cn"
	}
	cfg.TeslaCommandKeyPath = os.Getenv("TESLA_COMMAND_KEY_FILE")
	if cfg.TeslaCommandKeyPath == "" {
		cfg.TeslaCommandKeyPath = filepath.Join("public", ".well-known", "appspecific", "private-key.pem")
	}
	cfg.TeslaPartnerTokenURL = os.Getenv("TESLA_PARTNER_TOKEN_URL")
	if cfg.TeslaPartnerTokenURL == "" {
		cfg.TeslaPartnerTokenURL = "https://auth.tesla.cn/oauth2/v3/token"
	}
	cfg.TeslaPartnerScope = os.Getenv("TESLA_PARTNER_SCOPE")
	if cfg.TeslaPartnerScope == "" {
		cfg.TeslaPartnerScope = "openid user_data vehicle_device_data vehicle_cmds vehicle_charging_cmds vehicle_location offline_access"
	}
	cfg.TeslaPartnerDomain = os.Getenv("TESLA_PARTNER_DOMAIN")
	if cfg.TeslaPartnerDomain == "" {
		cfg.TeslaPartnerDomain = "dwdacbj25q.ap-southeast-1.awsapprunner.com"
	}

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
	cfg.DB.TimeZone = os.Getenv("DB_TIMEZONE")
	if cfg.DB.TimeZone == "" {
		cfg.DB.TimeZone = "UTC"
	}
	return cfg, nil
}

// loadEnv loads environment variables from a .env file. loadEnv 会从 .env 文件加载环境变量。
// It checks the current working directory first and then walks up the parent directories. 它会先检查当前工作目录，然后逐级向上查找父级目录。
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
