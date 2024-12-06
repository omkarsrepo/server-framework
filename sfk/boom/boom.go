// Unpublished Work © 2024

package boom

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Exception interface {
	Error() string
	IsEmpty() bool
	SetTraceId(string)
	SetValidationError(string)
	StatusCode() int
	responseObject() *exception
}

type exception struct {
	ErrorStatusCode int       `json:"statusCode,omitempty"`
	ErrorText       string    `json:"error,omitempty"`
	Message         string    `json:"message,omitempty"`
	ValidationError string    `json:"validationError,omitempty"`
	TraceId         string    `json:"traceId,omitempty"`
	Timestamp       time.Time `json:"timestamp,omitempty"`
}

func Boom(statusCode int, message string) Exception {
	return &exception{
		ErrorStatusCode: statusCode,
		ErrorText:       http.StatusText(statusCode),
		Message:         message,
		Timestamp:       time.Now(),
	}
}

func (e *exception) Error() string {
	return e.Message
}

func (e *exception) IsEmpty() bool {
	return e.Message == "" || e.ErrorStatusCode == 0
}

func (e *exception) SetTraceId(traceId string) {
	e.TraceId = traceId
}

func (e *exception) SetValidationError(validationError string) {
	e.ValidationError = validationError
}

func (e *exception) StatusCode() int {
	return e.ErrorStatusCode
}

func (e *exception) responseObject() *exception {
	return e
}

func (e *exception) Exception() Exception {
	return Boom(e.ErrorStatusCode, e.Message)
}

func InternalServerError() Exception {
	return Boom(http.StatusInternalServerError, "Something went wrong. Please retry in sometime or contact support team")
}

func PreconditionFailed(message string) Exception {
	return Boom(http.StatusPreconditionFailed, message)
}

func NotFound(message string) Exception {
	return Boom(http.StatusNotFound, message)
}

func Forbidden(message string) Exception {
	return Boom(http.StatusForbidden, message)
}

func Unauthorized(message string) Exception {
	return Boom(http.StatusUnauthorized, message)
}

func MethodNotAllowed(message string) Exception {
	return Boom(http.StatusMethodNotAllowed, message)
}

func BadRequest(message string) Exception {
	return Boom(http.StatusBadRequest, message)
}

func Abort(ginCtx *gin.Context, err Exception) {
	err.SetTraceId(ginCtx.GetString("TRACE_ID"))

	ginCtx.AbortWithStatusJSON(err.StatusCode(), err.responseObject())
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
	exp.SetTraceId(traceId)

	if validationError := getValidationError(err); validationError != "" {
		exp.SetValidationError(validationError)
	}

	ginCtx.AbortWithStatusJSON(statusCode, exp.responseObject())
}

func AbortForValidation(ginCtx *gin.Context, err error) {
	AbortForValidationWithMsg(ginCtx, err, "Request body validation failed")
}
