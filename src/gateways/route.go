package gateways

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RouteUsers(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/stripe")

	api.Post("/custom_price", gateway.InputPrice)
	api.Post("/webhook", gateway.TestWebhook)
	// api.Get("/payment", gateway.GetPaymentStatus)
	api.Get("/ws", websocket.New(gateway.WebSocketHandler))
}

func RouteIP(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/ip")
	api.Get("/check_ip", gateway.GetIp)
}
