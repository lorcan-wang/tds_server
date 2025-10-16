package main

import (
	"log"
	"tds_server/internal/config"
	"tds_server/internal/data"
	"tds_server/internal/repository"
	"tds_server/internal/router"
	"tds_server/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// 初始化数据库
	if err := data.InitDB(cfg); err != nil {
		log.Fatalf("init db error: %v", err)
	}

	// 构造repository
	tokenRepo := repository.NewTokenRepo()

	partnerSvc, err := service.NewPartnerTokenService(cfg)
	if err != nil {
		log.Fatalf("init partner token service error: %v", err)
	}

	commandSvc, err := service.NewVehicleCommandService(cfg)
	if err != nil {
		log.Fatalf("init vehicle command service error: %v", err)
	}

	r := router.NewRouter(cfg, tokenRepo, partnerSvc, commandSvc)

	addr := cfg.Server.Address

	log.Printf("Starting server on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
