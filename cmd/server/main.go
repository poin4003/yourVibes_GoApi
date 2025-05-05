package main

import (
	_ "github.com/poin4003/yourVibes_GoApi/cmd/swag/docs"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/initialize"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

// @title API Documentation YourVibes backend
// @version 1.0.0
// @description This is a sample YourVibes backend server
// @termsOfService https://github.com/poin4003/yourVibes_GoApi

// @contact.name TEAM HKTP
// @contact.url https://github.com/poin4003/yourVibes_GoApi
// @contact.email pchuy4003@gmail.com

// @host yourvibes.duckdns.org:8080
// @BasePath /v1/2024
// @schema https

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @description Token with 'Bearer ' prefix

func main() {
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DocExpansion("none"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	configMode := os.Getenv("YOURVIBES_SERVER_CONFIG_FILE")
	if configMode == "" {
		configMode = "dev"
	}

	if configMode == "prod" {
		certFile := "/etc/ssl/certs/fullchain.pem"
		keyFile := "/etc/ssl/certs/privkey.pem"
		log.Printf("Starting server in production mode with TLS on port :8080")
		if err := r.RunTLS(":8080", certFile, keyFile); err != nil {
			log.Fatalf("Failed to start server with TLS: %v", err)
		}
	} else {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
