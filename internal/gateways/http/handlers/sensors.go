package handlers

import (
	"homework/internal/domain"
	"homework/internal/gateways/http/middleware"
	"homework/internal/gateways/http/models"
	"homework/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type SensorsHandler struct {
	uc *usecase.Sensor
}

func NewSensorsHandler(uc *usecase.Sensor) *SensorsHandler {
	return &SensorsHandler{uc: uc}
}

func (h *SensorsHandler) SetupRouterGroup(r *gin.Engine) {
	sensorGroup := r.Group(h.GetPath())
	{
		sensorGroup.Use()
		sensorGroup.OPTIONS("", h.sensorsOptions)
		sensorGroup.POST("", middleware.ContentTypeJSONValidator(), h.registerSensor)
		sensorGroup.GET("", middleware.AcceptJSONValidator(), h.getSensors)
		sensorGroup.HEAD("", middleware.AcceptJSONValidator(), h.headSensors)
	}
}

func (h *SensorsHandler) GetPath() string {
	return "/sensors"
}

func toSensorModel(sensor domain.Sensor) *models.Sensor {
	v := &models.Sensor{
		ID:           &sensor.ID,
		SerialNumber: &sensor.SerialNumber,
		Type:         (*string)(&sensor.Type),
		CurrentState: &sensor.CurrentState,
		Description:  &sensor.Description,
		IsActive:     &sensor.IsActive,
		RegisteredAt: (*strfmt.DateTime)(&sensor.RegisteredAt),
		LastActivity: (*strfmt.DateTime)(&sensor.LastActivity),
	}
	return v
}

func toSensorsModel(sensors []domain.Sensor) []*models.Sensor {
	var s []*models.Sensor
	for _, sensor := range sensors {
		s = append(s, toSensorModel(sensor))
	}
	return s
}

func (h *SensorsHandler) registerSensor(ctx *gin.Context) {
	v := &models.SensorToCreate{}
	if err := ctx.ShouldBindJSON(v); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Error in the JSON format of the request body"})
		return
	}

	if err := v.Validate(nil); err != nil {
		reason := err.Error()
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "JSON body validation error: " + reason})
		return
	}

	sensor := domain.Sensor{
		Description:  *v.Description,
		IsActive:     *v.IsActive,
		SerialNumber: *v.SerialNumber,
		Type:         domain.SensorType(*v.Type),
	}
	out, err := h.uc.RegisterSensor(ctx, &sensor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Internal server error: unable to register sensor"})
		return
	}

	ctx.JSON(http.StatusOK, toSensorModel(*out))
}

func (h *SensorsHandler) getSensorsModel(ctx *gin.Context) ([]*models.Sensor, bool) {
	out, err := h.uc.GetSensors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Internal server error: unable to retrieve sensors"})
		return nil, false
	}

	sensors := toSensorsModel(out)
	return sensors, true
}

func (h *SensorsHandler) getSensors(ctx *gin.Context) {
	if sensors, ok := h.getSensorsModel(ctx); ok {
		ctx.JSON(http.StatusOK, sensors)
	}
}

func (h *SensorsHandler) headSensors(ctx *gin.Context) {
	if sensors, ok := h.getSensorsModel(ctx); ok {
		WriteHeaders(ctx, sensors)
	}
}

func (h *SensorsHandler) GetAvailableMethods() []string {
	return []string{http.MethodPost, http.MethodOptions, http.MethodGet, http.MethodHead}
}

func (h *SensorsHandler) sensorsOptions(ctx *gin.Context) {
	ctx.Header("Allow", strings.Join(h.GetAvailableMethods(), ","))
	ctx.Status(http.StatusNoContent)
}
