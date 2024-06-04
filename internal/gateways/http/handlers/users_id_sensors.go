package handlers

import (
	"homework/internal/gateways/http/middleware"
	"homework/internal/gateways/http/models"
	"homework/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SensorOwnerHandler struct {
	uc *usecase.User
}

func NewSensorOwnerHandler(uc *usecase.User) *SensorOwnerHandler {
	return &SensorOwnerHandler{uc: uc}
}

func (h *SensorOwnerHandler) SetupRouterGroup(r *gin.Engine) {
	sensorOwnerGroup := r.Group(h.GetPath())
	{
		sensorOwnerGroup.OPTIONS("", h.usersSensorsOptions)
		sensorOwnerGroup.POST("",
			middleware.ContentTypeJSONValidator(),
			h.bindSensorToUser,
		)
		sensorOwnerGroup.GET("",
			middleware.AcceptJSONValidator(),
			h.getUserSensors,
		)
		sensorOwnerGroup.HEAD("",
			middleware.AcceptJSONValidator(),
			h.headUserSensors,
		)
	}
}

func (h *SensorOwnerHandler) GetAvailableMethods() []string {
	return []string{http.MethodPost, http.MethodOptions, http.MethodGet, http.MethodHead}
}

func (h *SensorOwnerHandler) GetPath() string {
	return "/users/:user_id/sensors"
}

func (h *SensorOwnerHandler) bindSensorToUser(ctx *gin.Context) {
	u := &models.UserIDParam{}
	if err := ctx.ShouldBindUri(u); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the URI parameters of the request"})
		return
	}

	if err := u.Validate(nil); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "URI parameters validation error: " + err.Error()})
		return
	}

	s := &models.SensorToUserBinding{}
	if err := ctx.ShouldBindJSON(s); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Error in the JSON format of the request body"})
		return
	}

	if err := s.Validate(nil); err != nil {
		reason := err.Error()
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "JSON body validation error: " + reason})
		return
	}

	err := h.uc.AttachSensorToUser(ctx, *u.UserID, *s.SensorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"reason": "Sensor or user not found"})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *SensorOwnerHandler) getSensorsModel(ctx *gin.Context) ([]*models.Sensor, bool) {
	v := &models.UserIDParam{}
	if err := ctx.ShouldBindUri(v); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the URI parameters of the request"})
		return nil, false
	}

	if err := v.Validate(nil); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "URI parameters validation error: " + err.Error()})
		return nil, false
	}

	sensors, err := h.uc.GetUserSensors(ctx, *v.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"reason": "Sensors for the user not found"})
		return nil, false
	}

	return toSensorsModel(sensors), true
}

func (h *SensorOwnerHandler) getUserSensors(ctx *gin.Context) {
	if sensors, ok := h.getSensorsModel(ctx); ok {
		ctx.JSON(http.StatusOK, sensors)
	}
}

func (h *SensorOwnerHandler) headUserSensors(ctx *gin.Context) {
	if sensors, ok := h.getSensorsModel(ctx); ok {
		WriteHeaders(ctx, sensors)
	}
}

func (h *SensorOwnerHandler) usersSensorsOptions(ctx *gin.Context) {
	ctx.Header("Allow", strings.Join(h.GetAvailableMethods(), ","))
	ctx.Status(http.StatusNoContent)
}
