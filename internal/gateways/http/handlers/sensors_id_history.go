package handlers

import (
	"homework/internal/domain"
	"homework/internal/gateways/http/middleware"
	"homework/internal/gateways/http/models"
	"homework/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SensorHistoryHandler struct {
	uc *usecase.Event
}

func NewSensorHistoryHandler(uc *usecase.Event) *SensorHistoryHandler {
	return &SensorHistoryHandler{
		uc: uc,
	}
}

func (h *SensorHistoryHandler) GetPath() string {
	return "/sensors/:sensor_id/history"
}

func (h *SensorHistoryHandler) GetAvailableMethods() []string {
	return []string{http.MethodOptions, http.MethodGet, http.MethodHead}
}

func (h *SensorHistoryHandler) SetupRouterGroup(r *gin.Engine) {
	sensorDetailGroup := r.Group(h.GetPath())
	{
		sensorDetailGroup.GET("",
			middleware.AcceptJSONValidator(),
			h.getSensorsHistory,
		)

		sensorDetailGroup.HEAD("",
			middleware.AcceptJSONValidator(),
			h.headSensorsHistory,
		)

		sensorDetailGroup.OPTIONS("", h.sensorsHistoryOptions)
	}
}

func toSensorStatus(event domain.Event) *models.SensorStatus {
	return &models.SensorStatus{
		Timestamp: event.Timestamp,
		Payload:   event.Payload,
	}
}

func (h *SensorHistoryHandler) getSensorStatuses(ctx *gin.Context) ([]*models.SensorStatus, bool) {
	t := &models.TimeFraneQuery{}
	if err := ctx.ShouldBindQuery(t); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the query parameters of the request"})
		return nil, false
	}

	if err := t.Validate(nil); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Query parameters validation error: " + err.Error()})
		return nil, false
	}

	s := &models.SensorIDParam{}

	if err := ctx.ShouldBindUri(s); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the URI parameters of the request"})
		return nil, false
	}

	if err := s.Validate(nil); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "URI parameters validation error: " + err.Error()})
		return nil, false
	}

	events, err := h.uc.GetEventsByTimeFrame(ctx, *s.SensorID, *t.Start, *t.End)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"reason": err.Error()})
		return nil, false
	}

	var statuses []*models.SensorStatus

	for _, event := range events {
		statuses = append(statuses, toSensorStatus(event))
	}

	return statuses, true
}

func (h *SensorHistoryHandler) getSensorsHistory(ctx *gin.Context) {
	if statuses, ok := h.getSensorStatuses(ctx); ok {
		ctx.JSON(http.StatusOK, statuses)
	}
}

func (h *SensorHistoryHandler) headSensorsHistory(ctx *gin.Context) {
	if statuses, ok := h.getSensorStatuses(ctx); ok {
		WriteHeaders(ctx, statuses)
	}
}

func (h *SensorHistoryHandler) sensorsHistoryOptions(ctx *gin.Context) {
	ctx.Header("Allow", strings.Join(h.GetAvailableMethods(), ","))
	ctx.Status(http.StatusNoContent)
}
