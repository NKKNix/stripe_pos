package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		TimeFormat: "15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}
