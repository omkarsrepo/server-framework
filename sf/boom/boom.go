package boom

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Exception struct {
	StatusCode      int       `json:"statusCode,omitempty"`
	ErrorText       string    `json:"error,omitempty"`
	Message         string    `json:"message,omitempty"`
	ValidationError string    `json:"validationError,omitempty"`
	TraceId         string    `json:"traceId,omitempty"`
	Timestamp       time.Time `json:"timestamp,omitempty"`
}

func Boom(statusCode int, message string) *Exception {
	return &Exception{
		StatusCode: statusCode,
		ErrorText:  http.StatusText(statusCode),
		Message:    message,
		Timestamp:  time.Now(),
	}
}

func (props *Exception) Error() string {
	return props.Message
}

func (props *Exception) IsEmpty() bool {
	return props.Message == "" || props.StatusCode == 0
}

func InternalServerError() *Exception {
	return Boom(http.StatusInternalServerError, "Something went wrong. Please retry in sometime or contact support team")
}

func PreconditionFailed(message string) *Exception {
	return Boom(http.StatusPreconditionFailed, message)
}

func NotFound(message string) *Exception {
	return Boom(http.StatusNotFound, message)
}

func Forbidden(message string) *Exception {
	return Boom(http.StatusForbidden, message)
}

func Unauthorized(message string) *Exception {
	return Boom(http.StatusUnauthorized, message)
}

func MethodNotAllowed(message string) *Exception {
	return Boom(http.StatusMethodNotAllowed, message)
}

func BadRequest(message string) *Exception {
	return Boom(http.StatusBadRequest, message)
}

func Abort(ginCtx *gin.Context, err *Exception) {
	err.TraceId = ginCtx.GetString("TRACE_ID")
	ginCtx.AbortWithStatusJSON(err.StatusCode, err)
}

func getValidationError(err error) string {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		fieldError := err.(validator.ValidationErrors)[0]
		errorsStr := fmt.Sprintf("Field '%s' failed validation with constraint of '%s'", fieldError.Field(), fieldError.Tag())

		return errorsStr
	}

	return ""
}

func AbortForValidationWithMsg(ginCtx *gin.Context, err error, message string) {
	statusCode := http.StatusBadRequest
	traceId := ginCtx.GetString("TRACE_ID")

	exp := BadRequest(message)
	exp.TraceId = traceId

	if validationError := getValidationError(err); validationError != "" {
		exp.ValidationError = validationError
	}

	ginCtx.AbortWithStatusJSON(statusCode, exp)
}

func AbortForValidation(ginCtx *gin.Context, err error) {
	AbortForValidationWithMsg(ginCtx, err, "Request body validation failed")
}
