package main // Package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorest-api/internal/config"
	"gorest-api/internal/handler"
	"gorest-api/internal/logs"
	"gorest-api/internal/repository"
	"gorest-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	err := logs.InitLogger()
	if err != nil {
		log.Fatalf("Logger error: %s", err.Error())
	}

	appCfg, err := config.Init("configs")
	if err != nil {
		logs.Log().Error(err.Error())
		os.Exit(1)
	}

	fiberConfig := config.FiberConfig(appCfg)

	var mongoCfg config.MongoConfig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mongoClient, err := config.MongoNewClient(ctx, cancel, &mongoCfg)
	if err != nil {
		logs.Log().Error(err.Error())
		os.Exit(1)
	}

	app := fiber.New(fiberConfig)
	app.Use(logger.New())

	mongoDB := mongoClient.Database(mongoCfg.Name)
	appRepository := repository.NewRepository(mongoDB)
	appService := service.NewService(appRepository)
	appHandler := handler.NewHandler(appService)
	appHandler.InitRoutesFiber(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		oscall := <-c
		logs.Log().Info("Gracefully shutting down... ")
		logs.Log().Infof("System call: %s", oscall)
		cancel()
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error when shutting down...")
		}
	}()

	if err := app.Listen(":" + appCfg.HTTP.Port); err != nil {
		log.Panic(err)
	}

	logs.Log().Info("Running cleanup tasks...")
}
