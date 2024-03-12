package app

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/xlsx-generator/internal/handler"
	"github.com/illenko/xlsx-generator/internal/server"
	"github.com/illenko/xlsx-generator/internal/service"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type App struct{}

func (a App) Run() {
	fx.New(
		fx.Provide(
			zap.NewExample,
			service.New,
			handler.New,
			server.New,
		),
		fx.Invoke(func(e *gin.Engine) {
			err := e.Run(":8080")
			if err != nil {
				return
			}
		}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
