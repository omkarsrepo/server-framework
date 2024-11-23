package json

import (
	"github.com/gin-gonic/gin"
	"github.com/omkarsrepo/server-framework/sfk/boom"
	"strings"
)

func ExtractAuthorization(ginCtx *gin.Context) (string, *boom.Exception) {
	authHeader := ginCtx.GetHeader("Authorization")
	if len(authHeader) == 0 || !strings.Contains(authHeader, "Bearer ") {
		return "", boom.Unauthorized("Authorization header is invalid or empty")
	}

	token := strings.TrimSpace(strings.Split(authHeader, "Bearer ")[1])

	return token, nil
}
