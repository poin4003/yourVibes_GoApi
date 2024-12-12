package main

import (
	_ "github.com/poin4003/yourVibes_GoApi/cmd/swag/docs"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/initialize"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API Documentation YourVibes backend
// @version 1.0.0
// @description This is a sample YourVibes backend server
// @termsOfService https://github.com/poin4003/yourVibes_GoApi

// @contact.name TEAM HKTP
// @contact.url https://github.com/poin4003/yourVibes_GoApi
// @contact.email pchuy4003@gmail.com

// @host yourvibesapi.duckdns.org
// @BasePath /v1/2024
// @schema https

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @description Token without 'Bearer ' prefix

func main() {
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
