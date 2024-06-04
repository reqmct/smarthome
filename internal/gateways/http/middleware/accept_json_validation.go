package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AcceptJSONValidator() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		contentType := ctx.GetHeader("Accept")
		if contentType != "application/json" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"reason": "Content type is not application/json"})
			return
		}
		ctx.Next()
	}
}
