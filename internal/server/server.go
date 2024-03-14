package server

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/xlsx-generator/docs"
	"github.com/illenko/xlsx-generator/internal/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(handler handler.XlsxHandler, healthHandler handler.HealthHandler) *gin.Engine {
	e := gin.Default()

	e.POST("/xlsx", handler.Generate)
	e.GET("/health/liveness", healthHandler.Liveness)
	e.GET("/health/readiness", healthHandler.Readiness)

	docs.SwaggerInfo.BasePath = "/"
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return e
}
