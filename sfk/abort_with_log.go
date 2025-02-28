// Unpublished Work © 2024

package sfk

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/omkarsrepo/server-framework/sfk/boom"
)

func logError(ginCtx *gin.Context, err error) {
	logger := LoggerServiceInstance()

	req := ginCtx.Request
	cleanHeaders := cleanRequestHeaders(req.Header.Clone())

	logger.Error(ginCtx).
		Any("requestUrl", req.URL.Path).
		Any("requestMethod", req.Method).
		Any("requestHost", req.Host).
		Any("requestRemoteAddress", req.RemoteAddr).
		Any("requestClientIp", ginCtx.ClientIP()).
		Any("requestQuery", req.URL.Query()).
		Str("requestBody", ginCtx.GetString("STRING_REQ_BODY")).
		Any("requestHeaders", cleanHeaders).
		Any("requestContentLength", req.ContentLength).
		Str("received", fmt.Sprintf("%s https://%s%s", req.Method, req.Host, req.RequestURI)).
		Any("exception", err).
		Msg(err.Error())
}

func Abort(ginCtx *gin.Context, err error) {
	logError(ginCtx, err)

	boom.Abort(ginCtx, err)
}

func AbortForValidation(ginCtx *gin.Context, err error, message ...string) {
	logError(ginCtx, boom.BadRequest(err.Error()))

	if len(message) != 0 {
		boom.AbortForValidationWithMsg(ginCtx, err, message[0])
		return
	}

	boom.AbortForValidation(ginCtx, err)
}
