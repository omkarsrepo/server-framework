package boom

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Error struct {
	StatusCode int
	Message    string
	Body       gin.H
}

func (props *Error) Error() string {
	return props.Message
}

func (props *Error) IsEmpty() bool {
	return props.Message == "" || props.StatusCode == 0
}

func New(statusCode int, message string) *Error {
	body := gin.H{
		"error":      http.StatusText(statusCode),
		"statusCode": statusCode,
		"message":    message,
		"timestamp":  time.Now(),
	}

	return &Error{statusCode, message, body}
}

func InternalServerError() *Error {
	return New(http.StatusInternalServerError, "Something went wrong. Please retry in sometime or contact support team")
}

func AbortWithError(ctx *gin.Context, err *Error) {
	ctx.AbortWithStatusJSON(err.StatusCode, err.Body)
}

func getValidationErrors(err error) (interface{}, bool) {
	var errorsStr []string
	var validationErrors validator.ValidationErrors
	var nilObj interface{}

	if errors.As(err, &validationErrors) {
		for _, fieldError := range err.(validator.ValidationErrors) {
			errorsStr = append(errorsStr, fmt.Sprintf("Field '%s' failed validation with tag '%s'", fieldError.Field(), fieldError.Tag()))
		}

		return errorsStr, true
	}

	return nilObj, false
}

func AbortWithValidationErrors(ginCtx *gin.Context, err error, message string) {
	statusCode := http.StatusBadRequest

	if validationErrors, ok := getValidationErrors(err); ok {
		body := gin.H{
			"error":            http.StatusText(statusCode),
			"statusCode":       statusCode,
			"message":          message,
			"validationErrors": validationErrors,
			"timestamp":        time.Now(),
		}

		ginCtx.AbortWithStatusJSON(statusCode, &body)
	} else {
		body := gin.H{
			"error":      http.StatusText(statusCode),
			"statusCode": statusCode,
			"message":    message,
			"timestamp":  time.Now(),
		}

		ginCtx.AbortWithStatusJSON(statusCode, &body)
	}
}
