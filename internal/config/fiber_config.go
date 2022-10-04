package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	return fiber.Config{
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	}
}
