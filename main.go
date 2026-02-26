package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/suhas-developer07/EdwinNova-Server/internals/application"
)

func main() {
	mongoURI := getEnv("MONGO_URI", "mongodb+srv://suhas:Fordmustang1969@suhas.cbbha.mongodb.net/EdwinNova")
	dbName := getEnv("MONGO_DB", "edwinnova")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(dbName)

	repo := application.NewRepository(db)
	svc := application.NewService(repo)
	handler := application.NewHandler(svc, uploadDir)

	e := echo.New()

	e.POST("/applications", handler.CreateApplication)
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	if err := e.Start(os.Getenv("HTTP_ADDR")); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
