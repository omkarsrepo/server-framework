// Unpublished Work © 2024

package sfk

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func applyTraceHeader() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Set("TRACE_ID", uuid.New().String())

		traceId := strings.TrimSpace(ginCtx.GetHeader("X-Trace-ID"))

		if len(traceId) != 0 {
			ginCtx.Set("TRACE_ID", traceId)
		}

		ginCtx.Next()
	}
}
