package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/suhas-developer07/EdwinNova-Server/internals/email"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mongo"

	"github.com/suhas-developer07/EdwinNova-Server/internals/application"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbName := getEnv("MONGO_DB", "edwinnova")
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")

	/*Database Initialization */
	client,err := mongo.InitMongo(mongo.Config{
		URI: os.Getenv("MONGO_URI"),
		MaxPoolSize: 50,
		MinPoolSize: 5,
		Timeout: 30*time.Second,
	})

	if err != nil{
		panic(err)
	}

	db := client.Database(dbName)

	defer mongo.DisconnectMongo()

	/* RabbitMq Initialization */
	RabbitMQ_URI := os.Getenv("RABBITMQ_URI")

	rabbitmqConn,err := rabbitmq.New(RabbitMQ_URI)

	if err != nil {
		log.Fatalln("RabbitMq connection failed:Error",err)
	}

	defer rabbitmqConn.Close()
	
	queueName := "email_queue"

	err = rabbitmqConn.DeclareQueue(queueName,true)
	if err != nil{
		log.Fatalf("failed to declare email queue:%v",err)
	}

	publisher := email.NewPublisher(rabbitmqConn, queueName)

	
	/* Internals */
	repo := application.NewRepository(db)
	svc := application.NewService(repo,publisher)
	handler := application.NewHandler(svc, uploadDir)

	e := echo.New()

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
