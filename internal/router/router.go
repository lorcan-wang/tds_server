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
		protected.GET("/list", handler.GetList(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag", handler.GetVehicle(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/vehicle_data", handler.GetVehicleData(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/states", handler.GetVehicleStates(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/vehicle_data_request/:state", handler.GetVehicleDataRequest(cfg, tokenRepo))

		// Driver & sharing endpoints
		protected.GET("/vehicles/:vehicle_tag/drivers", handler.GetVehicleDrivers(cfg, tokenRepo))
		protected.DELETE("/vehicles/:vehicle_tag/drivers", handler.DeleteVehicleDriver(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/invitations", handler.GetVehicleShareInvites(cfg, tokenRepo))
		protected.POST("/vehicles/:vehicle_tag/invitations", handler.CreateVehicleShareInvite(cfg, tokenRepo))
		protected.POST("/vehicles/:vehicle_tag/invitations/:invitation_id/revoke", handler.RevokeVehicleShareInvite(cfg, tokenRepo))
		protected.POST("/invitations/redeem", handler.RedeemShareInvite(cfg, tokenRepo))

		// Vehicle capability endpoints
		protected.GET("/vehicles/:vehicle_tag/mobile_enabled", handler.GetVehicleMobileEnabled(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/nearby_charging_sites", handler.GetNearbyChargingSites(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/recent_alerts", handler.GetRecentAlerts(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/release_notes", handler.GetReleaseNotes(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/service_data", handler.GetServiceData(cfg, tokenRepo))

		// Fleet telemetry endpoints
		protected.POST("/vehicles/fleet_telemetry_config", handler.CreateFleetTelemetryConfig(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/fleet_telemetry_config", handler.GetFleetTelemetryConfig(cfg, tokenRepo))
		protected.DELETE("/vehicles/:vehicle_tag/fleet_telemetry_config", handler.DeleteFleetTelemetryConfig(cfg, tokenRepo))
		protected.POST("/vehicles/fleet_telemetry_config_jws", handler.PostFleetTelemetryConfigJWS(cfg, tokenRepo))
		protected.GET("/vehicles/:vehicle_tag/fleet_telemetry_errors", handler.GetFleetTelemetryErrors(cfg, tokenRepo))

		// Subscription endpoints
		protected.GET("/vehicle_subscriptions", handler.GetVehicleSubscriptions(cfg, tokenRepo))
		protected.POST("/vehicle_subscriptions", handler.SetVehicleSubscriptions(cfg, tokenRepo))
		protected.GET("/subscriptions", handler.GetSubscriptions(cfg, tokenRepo))
		protected.POST("/subscriptions", handler.SetSubscriptions(cfg, tokenRepo))
		protected.GET("/dx/vehicles/subscriptions/eligibility", handler.GetEligibleSubscriptions(cfg, tokenRepo))
		protected.GET("/dx/vehicles/upgrades/eligibility", handler.GetEligibleUpgrades(cfg, tokenRepo))

		// Misc endpoints
		protected.POST("/vehicles/fleet_status", handler.PostFleetStatus(cfg, tokenRepo))
		protected.GET("/dx/vehicles/options", handler.GetVehicleOptions(cfg, tokenRepo))
		protected.GET("/dx/warranty/details", handler.GetWarrantyDetails(cfg, tokenRepo))
		protected.POST("/vehicles/:vehicle_tag/signed_command", handler.SendSignedVehicleCommand(cfg, tokenRepo))

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
