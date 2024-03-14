package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type HealthHandler interface {
	Liveness(c *gin.Context)
	Readiness(c *gin.Context)
}

type HealthHandlerImpl struct {
	log *slog.Logger
}

func NewHealth(log *slog.Logger) HealthHandler {
	return HealthHandlerImpl{
		log: log,
	}
}

type HealthStatus string

const OkStatus HealthStatus = "OK"

type HealthResponse struct {
	Status HealthStatus `json:"status"`
}

var okHealthResponse = HealthResponse{Status: OkStatus}

// Liveness
//
//	@Tags			health
//	@Produce		json
//
//	@Success		200 {object} handler.HealthResponse
//	@Router			/health/liveness [get]
func (h HealthHandlerImpl) Liveness(c *gin.Context) {
	h.log.Info("Handled liveness check request")
	c.JSON(http.StatusOK, okHealthResponse)
}

// Readiness
//
//	@Tags			health
//	@Produce		json
//
//	@Success		200 {object} handler.HealthResponse
//	@Router			/health/readiness [get]
func (h HealthHandlerImpl) Readiness(c *gin.Context) {
	h.log.Info("Handled readiness check request")
	c.JSON(http.StatusOK, okHealthResponse)
}
