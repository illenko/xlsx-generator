package server

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/xlsx-generator/internal/handler"
)

func New(handler handler.XlsxHandler) *gin.Engine {
	e := gin.Default()

	e.POST("/xlsx", handler.Generate)

	return e
}
