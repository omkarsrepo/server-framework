package sf

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func ApplyTraceHeader() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Set("TRACE_ID", "")

		traceId := strings.TrimSpace(ginCtx.GetHeader("X-Trace-ID"))

		if len(traceId) != 0 {
			ginCtx.Set("TRACE_ID", traceId)
		}

		ginCtx.Next()
	}
}
