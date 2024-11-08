package main

import (
	"flag"
	"log"
	"os"

	"github.com/ncfex/dcart-gateway/internal/gateway"
	"github.com/ncfex/dcart-gateway/internal/infrastructure/config"
	"github.com/ncfex/dcart-gateway/pkg/httputil/response"
)

func main() {
	configPath := flag.String("config", "internal/infrastructure/config/config.example.yaml", "path to config file")
	flag.Parse()

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	responder := response.NewHTTPResponder(logger)

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	srv, err := gateway.NewGateway(cfg, logger, responder)
	if err != nil {
		log.Fatal("Failed to load server configuration:", err)
	}

	logger.Printf("starting server on port: %s", cfg.Server.Port)
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
