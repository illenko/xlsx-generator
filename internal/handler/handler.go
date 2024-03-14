package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/illenko/xlsx-generator/internal/logger"
	"github.com/illenko/xlsx-generator/internal/model"
	"github.com/illenko/xlsx-generator/internal/service"
	"log/slog"
	"net/http"
)

type XlsxHandler interface {
	Generate(c *gin.Context)
}

type XlsxHandlerImpl struct {
	log     *slog.Logger
	service service.XlsxService
}

func New(log *slog.Logger, service service.XlsxService) XlsxHandler {
	return XlsxHandlerImpl{
		log:     log,
		service: service,
	}
}

func (h XlsxHandlerImpl) Generate(c *gin.Context) {
	ctx := logger.AppendCtx(context.Background(), slog.String("request_id", uuid.New().String()))

	var req model.XlsxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	res, err := h.service.Generate(ctx, req)
	if err != nil {
		return
	}

	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Disposition", "attachment; filename=file.xlsx")
	c.Data(http.StatusOK, "application/octet-stream", res)
}
