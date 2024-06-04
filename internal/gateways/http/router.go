package http

import (
	"homework/internal/gateways/http/handlers"
	"homework/internal/gateways/http/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func setupRouter(r *gin.Engine, cases UseCases, wsHandler *WebSocketHandler) {
	endpoints := []handlers.Handler{
		handlers.NewUsersHandler(cases.User),
		handlers.NewSensorsHandler(cases.Sensor),
		handlers.NewSensorHandler(cases.Sensor),
		handlers.NewEventsHandler(cases.Event),
		handlers.NewSensorOwnerHandler(cases.User),
		handlers.NewSensorHistoryHandler(cases.Event),
	}

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, e := range endpoints {
		for _, method := range methods {
			if !contains(e.GetAvailableMethods(), method) {
				r.Handle(
					method,
					e.GetPath(),
					func(ctx *gin.Context) {
						ctx.AbortWithStatus(http.StatusMethodNotAllowed)
					},
				)
			}
		}
		e.SetupRouterGroup(r)
	}

	r.GET("/sensors/:sensor_id/events",
		func(ctx *gin.Context) {
			v := &models.SensorIDParam{}
			if err := ctx.ShouldBindUri(v); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the URI parameters of the request"})
				return
			}

			if err := v.Validate(nil); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "URI parameters validation error: " + err.Error()})
				return
			}
			err := wsHandler.Handle(ctx, *v.SensorID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"reason": err.Error()})
			}
		},
	)
}
