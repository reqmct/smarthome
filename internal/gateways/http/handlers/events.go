package handlers

import (
	"homework/internal/domain"
	"homework/internal/gateways/http/middleware"
	"homework/internal/gateways/http/models"
	"homework/internal/usecase"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EventsHandler struct {
	uc *usecase.Event
}

func NewEventsHandler(uc *usecase.Event) *EventsHandler {
	return &EventsHandler{
		uc: uc,
	}
}

func (h *EventsHandler) SetupRouterGroup(r *gin.Engine) {
	eventsGroup := r.Group(h.GetPath())
	{
		eventsGroup.OPTIONS("", h.eventsOptions)
		eventsGroup.POST("", middleware.ContentTypeJSONValidator(), h.registerEvent)
	}
}

func (h *EventsHandler) GetAvailableMethods() []string {
	return []string{http.MethodPost, http.MethodOptions}
}

func (h *EventsHandler) GetPath() string {
	return "/events"
}

func toEventModel(event domain.Event) *models.SensorEvent {
	v := &models.SensorEvent{
		Payload:            &event.Payload,
		SensorSerialNumber: &event.SensorSerialNumber,
	}

	return v
}

func (h *EventsHandler) registerEvent(ctx *gin.Context) {
	v := &models.SensorEvent{}
	if err := ctx.ShouldBindJSON(v); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Error in the JSON format of the request body"})
		return
	}

	if err := v.Validate(nil); err != nil {
		reason := err.Error()
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "JSON body validation error: " + reason})
		return
	}

	event := domain.Event{
		SensorSerialNumber: *v.SensorSerialNumber,
		Payload:            *v.Payload,
	}
	event.Timestamp = time.Now()
	err := h.uc.ReceiveEvent(ctx, &event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Internal server error: unable to process event"})
		return
	}

	ctx.JSON(http.StatusCreated, toEventModel(event))
}

func (h *EventsHandler) eventsOptions(ctx *gin.Context) {
	ctx.Header("Allow", strings.Join(h.GetAvailableMethods(), ","))
	ctx.Status(http.StatusNoContent)
}
