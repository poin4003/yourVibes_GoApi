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

// @host 192.168.20.113:8080
// @BasePath /v1/2024
// @schema http

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @description Token with 'Bearer ' prefix

func main() {
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DocExpansion("none"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))
	r.Run(":8080")
}
