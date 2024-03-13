package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/illenko/xlsx-generator/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type XlsxHandler interface {
	Generate(c *gin.Context)
}

type XlsxHandlerImpl struct {
	log     *zap.Logger
	service service.XlsxService
}

func New(log *zap.Logger, service service.XlsxService) XlsxHandler {
	return XlsxHandlerImpl{
		log:     log,
		service: service,
	}
}

func (h XlsxHandlerImpl) Generate(c *gin.Context) {
	var req model.XlsxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	res, err := h.service.Generate(req)
	if err != nil {
		return
	}

	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Disposition", "attachment; filename=file.xlsx")
	c.Data(http.StatusOK, "application/octet-stream", res)
}
