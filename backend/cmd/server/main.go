package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/internal/config"
	"github.com/hatamotoyuki/ninjo/backend/internal/handler"
	"github.com/hatamotoyuki/ninjo/backend/internal/infra"
	"github.com/hatamotoyuki/ninjo/backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	// DB接続 + マイグレーション
	client, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run migration: %v", err)
	}
	log.Println("database migration completed")

	// ファサード: DataStore → Usecase の順に組み立て
	ds := infra.NewDataStore(client)
	uc := usecase.NewUsecase(usecase.UsecaseConfig{
		DS:        ds,
		JWTSecret: cfg.JWTSecret,
	})

	// Echo セットアップ
	e := echo.New()
	e.Validator = handler.NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	handler.RegisterRoutes(e, uc)

	log.Fatal(e.Start(":" + cfg.Port))
}
