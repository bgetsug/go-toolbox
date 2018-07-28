package errors

import (
	"encoding/json"
	"net/http"

	"github.com/bgetsug/go-toolbox/logging"
	"github.com/bgetsug/go-toolbox/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"gopkg.in/go-playground/validator.v9"
)

// A middleware that outputs errors in a standard format.
func ErrorResponder() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.IsAborted() {
			return
		}

		if len(c.Errors) == 0 {
			return
		}

		var errorCollection []interface{}

		for _, contextError := range c.Errors {
			if validationErrors, ok := contextError.Err.(validator.ValidationErrors); ok {
				errorCollection = handleValidationErrors(errorCollection, validationErrors)
			} else if _, ok := contextError.Err.(*json.UnmarshalTypeError); ok {
				errorCollection = append(errorCollection, gin.H{
					"message": "Failed to parse JSON.",
					"code":    JSONUnmarshalFailed,
				})
			} else if _, ok := contextError.Err.(*json.SyntaxError); ok {
				errorCollection = append(errorCollection, gin.H{
					"message": "There was an error in the JSON syntax.",
					"code":    JSONDecodingFailed,
				})
			} else if errWithCode, ok := contextError.Err.(Error); ok {
				errorCollection = append(errorCollection, gin.H{
					"message": errWithCode.message,
					"code":    errWithCode.code,
				})
			} else {
				errorCollection = append(errorCollection, gin.H{
					"message": "An unknown error occurred.",
					"code":    Unknown,
				})
			}
		}

		c.JSON(c.Writer.Status(), gin.H{
			logging.RequestID: c.MustGet(logging.RequestID).(string),
			"errors":          errorCollection,
		})
	}
}

func handleValidationErrors(errorCollection []interface{}, validationErrors validator.ValidationErrors) []interface{} {
	for _, validationError := range validationErrors {
		errorCollection = append(errorCollection, gin.H{
			"message": validationError.Translate(validation.Translator),
			"code":    ValidationFailed,
		})
	}

	return errorCollection
}

func RecoveryHandler(c *gin.Context, err interface{}) {
	AbortWithServerError(c, Wrap(errors.New(err), PanicRecovery, "Recovered from panic"))
}

// A convenience function that calls AbortWithError with a Bad Request (400) status.
func AbortWithBadRequest(c *gin.Context, err error) {
	c.AbortWithError(http.StatusBadRequest, err) // nolint: errcheck
}

func AbortWithValidationError(c *gin.Context, err error) {
	c.AbortWithError(http.StatusUnprocessableEntity, err) // nolint: errcheck
}

func AbortWithServerError(c *gin.Context, err error) {
	c.AbortWithError(http.StatusInternalServerError, err) // nolint: errcheck
}
