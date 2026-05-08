package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// PostgreSQLドライバ（database/sqlに登録するためimport）
	_ "github.com/lib/pq"

	"github.com/hatamotoyuki/ninjo/backend/ent"
	"github.com/hatamotoyuki/ninjo/backend/internal/config"
)

func main() {
	cfg := config.Load()

	// entクライアント作成（PostgreSQLに接続）
	client, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()

	// 自動マイグレーション実行（entスキーマ → DBテーブル作成）
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run migration: %v", err)
	}
	log.Println("database migration completed")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	log.Fatal(e.Start(":" + cfg.Port))
}
