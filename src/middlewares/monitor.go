package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func MonitorMiddleware(path string) fiber.Handler {
	name, _ := GetModuleName()
	config := monitor.Config{}
	config.Title = fmt.Sprintf("%v Monitor", name)
	monitorHandler := monitor.New(config)

	return func(c *fiber.Ctx) error {
		if path == "" {
			path = "/monitor"
		}
		if c.Path() == path {
			// เรียกใช้ monitor handler
			return monitorHandler(c)
		}
		// ไปยัง handler หรือ middleware ถัดไป
		return c.Next()
	}
}
