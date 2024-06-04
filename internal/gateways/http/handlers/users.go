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

type UsersHandler struct {
	uc *usecase.User
}

func NewUsersHandler(uc *usecase.User) *UsersHandler {
	return &UsersHandler{uc: uc}
}

func (h *UsersHandler) SetupRouterGroup(r *gin.Engine) {
	userGroup := r.Group(h.GetPath())
	{
		userGroup.OPTIONS("", h.usersOptions)
		userGroup.POST("", middleware.ContentTypeJSONValidator(), h.createUser)
	}
}

func (h *UsersHandler) GetAvailableMethods() []string {
	return []string{http.MethodPost, http.MethodOptions}
}

func (h *UsersHandler) GetPath() string {
	return "/users"
}

func userValidator(user domain.User) *models.User {
	v := &models.User{
		ID:   &user.ID,
		Name: &user.Name,
	}

	return v
}

func (h *UsersHandler) createUser(ctx *gin.Context) {
	v := &models.UserToCreate{}
	if err := ctx.ShouldBindJSON(v); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "Error in the JSON format of the request body"})
		return
	}

	if err := v.Validate(nil); err != nil {
		reason := err.Error()
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "JSON body validation error: " + reason})
		return
	}

	user := domain.User{Name: *v.Name}
	out, err := h.uc.RegisterUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "Internal server error: unable to register user"})
		return
	}

	ctx.JSON(http.StatusOK, userValidator(*out))
}

func (h *UsersHandler) usersOptions(ctx *gin.Context) {
	ctx.Header("Allow", strings.Join(h.GetAvailableMethods(), ","))
	ctx.Status(http.StatusNoContent)
}
