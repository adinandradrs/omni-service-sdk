package service

import (
	"os"
	"time"

	"github.com/adinandradrs/omni-service-sdk/pkg/domain"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log(prod bool) (z *zap.Logger) {
	var c zap.Config
	if prod {
		c = zap.NewProductionConfig()
		c.DisableStacktrace = true

	} else {
		c = zap.NewDevelopmentConfig()
		c.DisableStacktrace = false

	}
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	z, _ = c.Build()
	return z

}

func Env(logger *zap.Logger) *domain.TechnicalError {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			logger.Error("failed to load configuration from .env file", zap.Error(err))
			return &domain.TechnicalError{
				Exception: err.Error(),
			}
		}
	} else {
		logger.Info("no configuration from .env")
	}
	return nil
}

func Exception(msg string, err error, logger *zap.Logger) *domain.TechnicalError {
	e := &domain.TechnicalError{
		Exception: err.Error(),
		Occurred:  time.Now().Unix(),
		Ticket:    uuid.New().String(),
	}
	logger.Error(msg, zap.Any("", e))
	return e
}
