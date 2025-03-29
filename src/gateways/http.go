package gateways

import (
	"go-fiber-template/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	IPService      services.IIpService
	StripeService  services.IStripeService
	OrderService   services.IOrderService // Inject OrderService interface
}

// NewHTTPGateway - Initializes HTTP routes and injects services
func NewHTTPGateway(app *fiber.App, ip services.IIpService, stripe services.IStripeService, order services.IOrderService) {
	gateway := &HTTPGateway{
		IPService:    ip,
		StripeService: stripe,
		OrderService:  order,
	}

	RouteUsers(*gateway, app)
	RouteIP(*gateway, app)
}
