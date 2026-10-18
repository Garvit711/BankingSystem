package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Garvit711/bank/internal/platform"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		platform.RespondError(c, http.StatusBadRequest, "invalid input")
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			platform.RespondError(c, http.StatusConflict, "email already registered")
		case errors.Is(err, ErrPasswordTooLong):
			platform.RespondError(c, http.StatusBadRequest, "password too long")
		default:
			log.Printf("register failed: %v", err)
			platform.RespondError(c, http.StatusInternalServerError, "internal error")
		}
		return
	}

	platform.RespondJSON(c, http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		platform.RespondError(c, http.StatusBadRequest, "invalid input")
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			platform.RespondError(c, http.StatusUnauthorized, "invalid email or password")
			return
		}
		log.Printf("login failed: %v", err)
		platform.RespondError(c, http.StatusInternalServerError, "internal error")
		return
	}

	platform.RespondJSON(c, http.StatusOK, gin.H{"token": token})
}