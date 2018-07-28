package stats

import (
	"net/http"

	"github.com/dimiro1/health"
	"github.com/gin-gonic/gin"
)

var CompositeChecker health.CompositeChecker

func HealthChecks() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := CompositeChecker.Check()

		httpCode := http.StatusOK

		if status.IsDown() {
			httpCode = http.StatusServiceUnavailable
		}

		c.JSON(httpCode, status)
	}
}
