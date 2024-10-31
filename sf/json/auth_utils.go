package json

import (
	"github.com/gin-gonic/gin"
	"github.com/omkarsrepo/server-framework/sf/boom"
	"net/http"
	"strings"
)

func ExtractAuthorization(ginCtx *gin.Context) (string, *boom.Error) {
	authHeader := ginCtx.GetHeader("Authorization")
	if len(authHeader) == 0 || !strings.Contains(authHeader, "Bearer ") {
		return "", boom.New(http.StatusUnauthorized, "Authorization header is invalid or empty")
	}

	token := strings.TrimSpace(strings.Split(authHeader, "Bearer ")[1])

	return token, nil
}
