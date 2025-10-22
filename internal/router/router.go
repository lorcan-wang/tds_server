package router

import (
	"net/http"
	"path/filepath"
	"runtime"
	"tds_server/internal/config"
	"tds_server/internal/handler"
	"tds_server/internal/middleware"
	"tds_server/internal/repository"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, tokenRepo *repository.TokenRepo, partnerSvc *service.PartnerTokenService, commandSvc *service.VehicleCommandService) *gin.Engine {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ping"})
	})

	if partnerSvc != nil {
		r.Use(func(c *gin.Context) {
			c.Set("partnerTokenService", partnerSvc)
			c.Next()
		})
	}

	publicKeyFile := publicKeyFilePath()
	r.StaticFile("/.well-known/appspecific/com.tesla.3p.public-key.pem", publicKeyFile)

	api := r.Group("/api")
	{
		api.GET("/login", handler.LoginRedirect(cfg))
		api.GET("/login/callback", handler.LoginCallback(cfg, tokenRepo))

		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(cfg))
		protected.GET("/1/vehicles", handler.ListVehicles(cfg, tokenRepo))
		protected.GET("/1/vehicles/:vehicle_tag", handler.GetVehicle(cfg, tokenRepo))
		protected.GET("/1/vehicles/:vehicle_tag/vehicle_data", handler.GetVehicleData(cfg, tokenRepo))
		protected.GET("/1/vehicles/:vehicle_tag/drivers", handler.GetVehicleDrivers(cfg, tokenRepo))
		protected.POST("/vehicles/:vehicle_tag/command/*command_path", handler.VehicleCommand(cfg, tokenRepo, commandSvc))
	}
	return r
}

func publicKeyFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("public", ".well-known", "appspecific", "com.tesla.3p.public-key.pem")
	}

	baseDir := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(baseDir, "public", ".well-known", "appspecific", "com.tesla.3p.public-key.pem")
}
