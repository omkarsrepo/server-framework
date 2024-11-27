package sfk

import (
	"bytes"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/omkarsrepo/server-framework/sfk/json"
	"io"
	"net/http"
)

type requestLoggerMiddleware struct {
	logger LoggerService
}

func purgeField(decodedBody *map[string]interface{}, fieldName string) {
	if _, err := json.AnyValueOf[any](decodedBody, fieldName); err == nil {
		delete(*decodedBody, fieldName)
	}
}

func cleanRequestBody(body io.ReadCloser, ctx *gin.Context) map[string]interface{} {
	defer body.Close()
	var decodedBody map[string]interface{}

	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err == nil && jsoniter.Unmarshal(bodyBytes, &decodedBody) == nil {
			purgeField(&decodedBody, "password")
			purgeField(&decodedBody, "clientSecret")
			purgeField(&decodedBody, "sessionToken")
		}
	}

	return decodedBody
}

func cleanRequestHeaders(headers http.Header) http.Header {
	if headers.Get("Authorization") != "" {
		headers.Del("Authorization")
	}
	return headers
}

func (r *requestLoggerMiddleware) applyFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		cleanBody := cleanRequestBody(req.Body, ctx)
		cleanHeaders := cleanRequestHeaders(req.Header.Clone())

		r.logger.Info(ctx).
			Any("requestUrl", req.URL.Path).
			Any("requestMethod", req.Method).
			Any("requestHost", req.Host).
			Any("requestRemoteAddress", req.RemoteAddr).
			Any("requestClientIp", ctx.ClientIP()).
			Any("requestQuery", req.URL.Query()).
			Any("requestBody", cleanBody).
			Any("requestHeaders", cleanHeaders).
			Any("requestContentLength", req.ContentLength).
			Msgf("received %s https://%s%s", req.Method, req.Host, req.RequestURI)

		ctx.Next()
	}
}

func applyRequestLoggerMiddleware() gin.HandlerFunc {
	return (&requestLoggerMiddleware{logger: LoggerServiceInstance()}).applyFilter()
}
