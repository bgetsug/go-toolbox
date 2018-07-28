package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, accounts gin.Accounts, appConfig interface{}) {
	router.LoadHTMLFiles("vendor/github.com/bgetsug/go-toolbox/stats/buildmetadata.tmpl")

	statsRoutes := router.Group("/stats")
	{
		statsRoutes.GET("/health", HealthChecks())
		statsRoutes.GET("/build", BuildMetadata(router))
		statsRoutes.GET("/config", gin.BasicAuth(accounts), func(c *gin.Context) {
			c.JSON(http.StatusOK, appConfig)
		})
	}
}
