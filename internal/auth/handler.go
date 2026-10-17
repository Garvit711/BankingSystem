package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/Garvit711/bank/internal/platform"
)

func Register(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. read and validate the request
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			platform.RespondError(c, http.StatusBadRequest, "invalid input")
			return
		}
		email := strings.ToLower(strings.TrimSpace(req.Email))

		// 2. hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			platform.RespondError(c, http.StatusBadRequest, "invalid password")
			return
		}

		// 3. insert the user
		user := User{Email: email}
		err = db.QueryRow(c.Request.Context(),
			`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`,
			email, string(hash),
		).Scan(&user.ID, &user.CreatedAt)

		// 4. handle the result
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				platform.RespondError(c, http.StatusConflict, "email already registered")
				return
			}
			log.Printf("register failed: %v", err)
			platform.RespondError(c, http.StatusInternalServerError, "internal error")
			return
		}

		platform.RespondJSON(c, http.StatusCreated, user)
	}
}