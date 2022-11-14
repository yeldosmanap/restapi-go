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

	"gorestapi/internal/config"
	"gorestapi/internal/handler"
	"gorestapi/internal/logs"
	"gorestapi/internal/middlewares"
	"gorestapi/internal/repository"
	"gorestapi/internal/service"

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

	prometheusInstance := middlewares.New("restapi-go")
	prometheusInstance.RegisterAt(app, "/api/v1/metrics")

	app.Use(prometheusInstance.Middleware)

	appHandler.InitRoutesFiber(app)

	go start(app, appCfg.HTTP.Port)

	stopChannel, closeChannel := createChannel()
	defer closeChannel()

	logs.Log().Info("Notified ", <-stopChannel)
	shutdown(ctx, app)
}

func start(server *fiber.App, port string) {
	logs.Log().Info("Application started")
	if err := server.Listen(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		logs.Log().Info("application stopped gracefully")
	}
}

func createChannel() (chan os.Signal, func()) {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopChannel, func() {
		close(stopChannel)
	}
}

func shutdown(ctx context.Context, app *fiber.App) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		panic(err)
	} else {
		logs.Log().Info("Application shutdown")
	}
}
