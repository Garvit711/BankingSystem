package account

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Garvit711/bank/internal/auth"
	"github.com/Garvit711/bank/internal/platform"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	acc, err := h.svc.Create(c.Request.Context(), auth.UserID(c))
	if err != nil {
		log.Printf("create account failed: %v", err)
		platform.RespondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	platform.RespondJSON(c, http.StatusCreated, acc)
}

func (h *Handler) List(c *gin.Context) {
	accounts, err := h.svc.List(c.Request.Context(), auth.UserID(c))
	if err != nil {
		log.Printf("list accounts failed: %v", err)
		platform.RespondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	platform.RespondJSON(c, http.StatusOK, accounts)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		platform.RespondError(c, http.StatusBadRequest, "invalid account id")
		return
	}

	acc, err := h.svc.Get(c.Request.Context(), id, auth.UserID(c))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			platform.RespondError(c, http.StatusNotFound, "account not found")
			return
		}
		log.Printf("get account failed: %v", err)
		platform.RespondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	platform.RespondJSON(c, http.StatusOK, acc)
}
