package auth

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/Garvit711/bank/internal/platform"
)

func Register(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			platform.RespondError(c, http.StatusBadRequest, "invalid input")
			return
		}
		email := strings.ToLower(strings.TrimSpace(req.Email))

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			platform.RespondError(c, http.StatusBadRequest, "invalid password")
			return
		}
		user := User{Email: email}
		err = db.QueryRow(c.Request.Context(),
			`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`,
			email, string(hash),
		).Scan(&user.ID, &user.CreatedAt)

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

func Login(db *pgxpool.Pool, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			platform.RespondError(c, http.StatusBadRequest, "invalid input")
			return
		}
		email := strings.ToLower(strings.TrimSpace(req.Email))
		var user User
		err := db.QueryRow(c.Request.Context(), `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				platform.RespondError(c, http.StatusUnauthorized, "invalid email or password")
				return
			}
			log.Printf("login query failed: %v", err)
			platform.RespondError(c, http.StatusInternalServerError, "internal error")
			return
		}
		err1 := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err1 != nil {
			platform.RespondError(c, http.StatusUnauthorized, "Invalid Email or Password")
			return

		}
		claims := jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
		if err != nil {
			log.Printf("Invalid Token Generation %v", err)
			platform.RespondError(c, http.StatusInternalServerError, "internal error")
			return

		}
		platform.RespondJSON(c, http.StatusOK, gin.H{"token": token})
	}
}
