package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

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

	r := gin.Default()
	r.POST("/auth/register", auth.Register(pool))
	r.POST("/auth/login", auth.Login(pool, cfg.JWTSecret))

	log.Fatal(r.Run(":" + cfg.Port))
}