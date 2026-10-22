package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Garvit711/bank/internal/account"
	"github.com/Garvit711/bank/internal/auth"
	"github.com/Garvit711/bank/internal/platform"

)

func main() {
	cfg, err := platform.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	pool, err := platform.NewDB(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer pool.Close()
	log.Println("connected to db")
	repo := auth.NewRepository(pool)
	svc := auth.NewService(repo, cfg.JWTSecret)
	h := auth.NewHandler(svc)

	r := gin.Default()
	r.GET("/me", auth.RequireAuth(cfg.JWTSecret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": auth.UserID(c)})
	})
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}

	// Account Handling 
	accountRepo := account.NewRepository(pool)
	accountSvc := account.NewService(accountRepo)
	accountHandler := account.NewHandler(accountSvc)
	protected := r.Group("/", auth.RequireAuth(cfg.JWTSecret))
{
	protected.POST("/accounts", accountHandler.Create)
	protected.GET("/accounts", accountHandler.List)
	protected.GET("/accounts/:id", accountHandler.Get)
}



if os.Getenv("ENABLE_DEV_ROUTES") == "true" {
	protected.POST("/accounts/:id/deposit", accountHandler.Deposit)
}

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Printf("listening on :%s", cfg.Port)
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-sigCtx.Done()

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
