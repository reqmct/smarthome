package handlers

import (
	"homework/internal/gateways/http/middleware"
	"homework/internal/gateways/http/models"
	"homework/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	uc *usecase.Sensor
}

func NewSensorHandler(uc *usecase.Sensor) *SensorHandler {
	return &SensorHandler{uc: uc}
}

func (h *SensorHandler) SetupRouterGroup(r *gin.Engine) {
	sensorDetailGroup := r.Group(h.GetPath())
	{
		sensorDetailGroup.OPTIONS("", h.sensorOptions)
		sensorDetailGroup.GET("",
			middleware.AcceptJSONValidator(),
			h.getSensor,
		)
		sensorDetailGroup.HEAD("",
			middleware.AcceptJSONValidator(),
			h.headSensor,
		)
	}
}

func (h *SensorHandler) getSensorModel(ctx *gin.Context) (*models.Sensor, bool) {
	v := &models.SensorIDParam{}
	if err := ctx.ShouldBindUri(v); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the URI parameters of the request"})
		return nil, false
	}

	if err := v.Validate(nil); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "URI parameters validation error: " + err.Error()})
		return nil, false
	}
	sensor, err := h.uc.GetSensorByID(ctx, *v.SensorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"reason": "Sensor not found"})
		return nil, false
	}

	return toSensorModel(*sensor), true
}

func (h *SensorHandler) getSensor(ctx *gin.Context) {
	if sensor, ok := h.getSensorModel(ctx); ok {
		ctx.JSON(http.StatusOK, sensor)
	}
}

func (h *SensorHandler) headSensor(ctx *gin.Context) {
	if sensor, ok := h.getSensorModel(ctx); ok {
		WriteHeaders(ctx, sensor)
	}
}

func (h *SensorHandler) sensorOptions(ctx *gin.Context) {
	ctx.Header("Allow", strings.Join(h.GetAvailableMethods(), ","))
	ctx.Status(http.StatusNoContent)
}

func (h *SensorHandler) GetPath() string {
	return "/sensors/:sensor_id"
}

func (h *SensorHandler) GetAvailableMethods() []string {
	return []string{http.MethodOptions, http.MethodGet, http.MethodHead}
}
