package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetPath() string
	GetAvailableMethods() []string
	SetupRouterGroup(r *gin.Engine)
}

func WriteHeaders(ctx *gin.Context, source any) {
	data, err := json.Marshal(source)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Internal server error: unable to marshal sensor data"})
		return
	}

	ctx.Header("Content-Length", strconv.Itoa(len(data)))

	ctx.Status(http.StatusOK)
}
