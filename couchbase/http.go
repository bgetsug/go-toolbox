package couchbase

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// Register routes for bootstrapping databases
func RegisterCbBootstrapRoutes(router *gin.Engine, accounts gin.Accounts) {
	dbRoutes := router.Group("/boot/cb", gin.BasicAuth(accounts))
	{
		dbRoutes.PUT("/indexes", func(c *gin.Context) {
			var indexCollection []interface{}
			var indexErrors []interface{}

			indexes, errors := Cb.CreateIndexes()
			indexCollection = append(indexCollection, indexes)

			for _, err := range errors {
				indexErrors = append(indexErrors, err)
			}

			status := http.StatusCreated

			if len(indexErrors) > 0 {
				status = http.StatusInternalServerError
			}

			c.JSON(status, gin.H{"indexes": indexCollection, "error": indexErrors})
		})

		dbRoutes.PUT("/seeds", func(c *gin.Context) {
			seederResults := make(chan SeederResults)

			go Cb.Seed(seederResults)

			type seedContainer struct {
				Type string      `json:"type"`
				Data interface{} `json:"data"`
			}

			var seedCollection []seedContainer
			var errorCollection []string

			for results := range seederResults {

				for _, seed := range results.Seeds {
					seedType := "UNKNOWN"

					if structs.IsStruct(seed) {
						seedType = structs.Name(seed)
					}

					seedCollection = append(seedCollection, seedContainer{Type: seedType, Data: seed})
				}

				for _, err := range results.Errors {
					errorCollection = append(errorCollection, err.Error())
				}
			}

			status := http.StatusCreated

			if len(errorCollection) > 0 {
				status = http.StatusInternalServerError
			}

			c.JSON(status, gin.H{"seeds": seedCollection, "errors": errorCollection})
		})
	}
}
