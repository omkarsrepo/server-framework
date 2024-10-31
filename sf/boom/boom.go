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
	StatusCode       int         `json:"statusCode,omitempty"`
	Exception        string      `json:"error,omitempty"`
	Message          string      `json:"message,omitempty"`
	ValidationErrors interface{} `json:"validationErrors,omitempty"`
	TraceId          string      `json:"traceId,omitempty"`
	Timestamp        time.Time   `json:"timestamp,omitempty"`
}

func (props *Error) Error() string {
	return props.Message
}

func (props *Error) IsEmpty() bool {
	return props.Message == "" || props.StatusCode == 0
}

func New(statusCode int, message string) *Error {
	return &Error{
		StatusCode:       statusCode,
		Exception:        http.StatusText(statusCode),
		Message:          message,
		Timestamp:        time.Now(),
		ValidationErrors: nil,
	}
}

func InternalServerError() *Error {
	return New(http.StatusInternalServerError, "Something went wrong. Please retry in sometime or contact support team")
}

func AbortWithError(ctx *gin.Context, err *Error) {
	err.TraceId = ctx.GetString("TRACE_ID")
	ctx.AbortWithStatusJSON(err.StatusCode, err)
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
	traceId := ginCtx.GetString("TRACE_ID")

	exp := New(statusCode, message)
	exp.TraceId = traceId

	if validationErrors, ok := getValidationErrors(err); ok {
		exp.ValidationErrors = validationErrors
	}

	ginCtx.AbortWithStatusJSON(statusCode, exp)
}
