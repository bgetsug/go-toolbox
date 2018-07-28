package logging

import (
	"time"

	"github.com/bgetsug/go-toolbox/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

const (
	RequestID = "requestId"
)

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		requestID := xid.New().String()
		c.Set(RequestID, requestID)

		start := time.Now()
		c.Next()
		end := time.Now()

		latency := end.Sub(start)

		log := logger.With(
			zap.String("module", "http"),
			zap.String(RequestID, requestID),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.UTC().Format(time.RFC3339)),
			zap.Any("headers", c.Request.Header),
			zap.Any("params", c.Params),
		)

		if len(c.Errors) > 0 {
			if c.Writer.Status() >= 500 {
				for _, contextError := range c.Errors {
					withContextError(log, contextError).Error(contextError.Error())
				}
			} else if c.Writer.Status() >= 400 {
				for _, contextError := range c.Errors {
					validationErrors, isValidationErrors := contextError.Err.(validator.ValidationErrors)

					if isValidationErrors {
						log := *log

						for _, validationError := range validationErrors {
							log.Warn(validationError.Translate(validation.Translator))
						}

						return
					}

					withContextError(log, contextError).Warn(contextError.Error())
				}
			}

			return
		}

		log.Info("")
	}
}

func withContextError(logger *zap.Logger, contextError *gin.Error) *zap.Logger {
	return logger.With(zap.Any("gin.Error.Meta", contextError.Meta), zap.Error(contextError))
}
