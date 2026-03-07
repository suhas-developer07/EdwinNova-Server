package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/suhas-developer07/EdwinNova-Server/internals/application"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mail"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mongo"
	storage "github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/s3"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbName := getEnv("MONGO_DB", "edwinnova")

	/* MongoDB */
	client, err := mongo.InitMongo(mongo.Config{
		URI:         os.Getenv("MONGO_URI"),
		MaxPoolSize: 50,
		MinPoolSize: 5,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)
	defer mongo.DisconnectMongo()

	/* SMTP */
	smtpClient, err := mail.NewSMTPClient()
	if err != nil {
		log.Fatalln("Failed to initialize SMTP client:", err)
	}

	/* AWS Config */
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("AWS region:", cfg.Region)

	s3Client := awss3.NewFromConfig(cfg)

	s3Storage := storage.NewS3Storage(
		s3Client,
		os.Getenv("FILES_BUCKET"),
	)

	/* Internals */
	repo := application.NewRepository(db)
	svc := application.NewService(repo, smtpClient)

	handler := application.NewHandler(svc, s3Storage)

	/* Echo */
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))

	e.POST("/applications", handler.CreateApplication)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
