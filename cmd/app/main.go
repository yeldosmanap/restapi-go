package main // Package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"gorest-api/internal/config"
	"gorest-api/internal/handler"
	"gorest-api/internal/logs"
	"gorest-api/internal/repository"
	"gorest-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// @title           Go REST API
// @version         1.0
// @description Projects REST API for Go.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	_ = logs.InitLogger()

	fiberConfig := config.FiberConfig()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	cfg, err := config.Init("configs")
	if err != nil {
		logs.Log().Error(err.Error())
	}

	mongoClient, err := config.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logs.Log().Error(err.Error())
		return
	}

	app := fiber.New(fiberConfig)
	app.Use(logger.New())

	mongo := mongoClient.Database(cfg.Mongo.Name)
	appRepository := repository.NewRepository(mongo)
	appService := service.NewService(appRepository)
	appHandler := handler.NewHandler(appService)
	appHandler.InitRoutesFiber(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}

	fmt.Println("Running cleanup tasks...")
}
