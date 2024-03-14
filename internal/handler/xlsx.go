package handler

import (
	"bufio"
	"bytes"
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

func NewXlsx(log *slog.Logger, service service.XlsxService) XlsxHandler {
	return XlsxHandlerImpl{
		log:     log,
		service: service,
	}
}

type ResponseError struct {
	Id    uuid.UUID `json:"id"`
	Error string    `json:"error"`
}

// Generate
//
//	@Summary	Generates xlsx file
//	@Tags			xlsx
//	@Accept			json
//	@Produce		octet-stream
//
// @Param data body model.XlsxRequest true "Request"
//
//	@Success		200
//	@Failure		400 {object} handler.ResponseError
//	@Failure		500 {object} handler.ResponseError
//	@Router			/xlsx [post]
func (h XlsxHandlerImpl) Generate(c *gin.Context) {
	requestId := uuid.New()
	ctx := logger.AppendCtx(context.Background(), slog.String("requestId", requestId.String()))

	h.log.InfoContext(ctx, "Processing generation request")

	var req model.XlsxRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Sheets == nil {
		h.log.WarnContext(ctx, "Invalid request body")
		c.JSON(http.StatusBadRequest, ResponseError{Id: requestId, Error: "Invalid Request Body"})
		return
	}

	wb := h.service.Generate(ctx, req)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err := wb.Write(writer)

	if err != nil {
		h.log.ErrorContext(ctx, "Internal server error")
		c.JSON(http.StatusInternalServerError, ResponseError{Id: requestId, Error: "Internal Server Error"})
		return
	}

	h.log.InfoContext(ctx, "Successfully generated XLSX file")

	c.Header("Content-Disposition", "attachment; filename=file.xlsx")
	c.Data(http.StatusOK, "application/vnd.ms-excel", b.Bytes())
}
