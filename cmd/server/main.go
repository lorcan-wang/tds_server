package main

import (
	"log"
	"tds_server/internal/config"
	"tds_server/internal/data"
	"tds_server/internal/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := data.InitDB(cfg); err != nil {
		log.Fatalf("init db error: %v", err)
	}

	r := router.NewRouter(cfg)

	addr := cfg.Server.Address

	log.Printf("Starting server on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
