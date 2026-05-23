package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/inherit-dapp/chain/backend/api"
	"github.com/inherit-dapp/chain/backend/config"
	"github.com/inherit-dapp/chain/backend/db"
	"github.com/inherit-dapp/chain/backend/monitor"
)

func main() {
	log.Println("=== Legacy DApp Backend Service ===")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: server=%s:%d, db=%s:%d",
		cfg.ServerHost, cfg.ServerPort, cfg.DBHost, cfg.DBPort)

	// Connect to database
	database, err := db.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Running without database (some features will be unavailable)")
		database = nil
	}
	if database != nil {
		defer database.Close()
		log.Println("Database connected successfully")
	}

	// Start heartbeat monitor
	if cfg.MonitorEnabled && database != nil {
		heartbeatMonitor := monitor.NewHeartbeatMonitor(cfg, database)
		go heartbeatMonitor.Start()
		defer heartbeatMonitor.Stop()

		keeperBot := monitor.NewKeeperBot(cfg, database)
		go keeperBot.Start()
		defer keeperBot.Stop()

		log.Println("Heartbeat monitor and keeper bot started")
	}

	// Start API server
	server := api.NewServer(cfg, database)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down...", sig)
		os.Exit(0)
	}()

	log.Printf("Starting API server on %s:%d", cfg.ServerHost, cfg.ServerPort)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
