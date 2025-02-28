package sfk

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func applyStringReqBody() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Set("STRING_REQ_BODY", "")

		reqBody := ginCtx.Request.Body
		defer reqBody.Close()

		if reqBody != nil {
			body, err := io.ReadAll(reqBody)
			if err == nil {
				ginCtx.Set("STRING_REQ_BODY", string(body))
			}

			ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		ginCtx.Next()
	}
}
