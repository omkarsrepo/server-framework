package sf

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/omkarsrepo/server-framework/sf/boom"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func ApplyRequestTimeout() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		ginCtx.Request = ginCtx.Request.WithContext(ctx)

		ginCtx.Next()

		if ctx.Err() != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
			exp := boom.New(http.StatusGatewayTimeout, "Request took too long. Please retry after sometime or contact support!")
			ginCtx.AbortWithStatusJSON(exp.StatusCode, &exp.Body)
		}
	}
}
