package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ContentTypeJSONValidator() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		contentType := ctx.GetHeader("Content-Type")
		if contentType != "application/json" {
			ctx.JSON(http.StatusUnsupportedMediaType, gin.H{"reason": "Content type not supported, expected application/json"})
			return
		}
		ctx.Next()
	}
}
