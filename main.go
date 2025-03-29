package main

import (
	"fmt"
	"go-fiber-template/src/configuration"
	"go-fiber-template/src/domain/repositories"
	"go-fiber-template/src/gateways"
	"go-fiber-template/src/infrastructure/httpclients"
	"go-fiber-template/src/middlewares"
	"go-fiber-template/src/services"
	sv "go-fiber-template/src/services"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	supa "github.com/supabase-community/supabase-go"
)

var Supabase *supa.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	app := fiber.New(configuration.NewFiberConfiguration())

	app.Use(middlewares.ScalarMiddleware(middlewares.Config{
		PathURL: "/api/docs",
		SpecURL: "./src/docs/swagger.yaml",
	}))
	app.Use(middlewares.MonitorMiddleware("/api/monitor"))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middlewares.Logger())

	supabaseURL := os.Getenv("VITE_SUPABASE_URL")
	supabaseKey := os.Getenv("VITE_SUPABASE_KEY")
	var err error
	Supabase, err = supa.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		fmt.Println("Failed to create Supabase client: ", err)
	}
	orderRepo := repositories.NewOrderRepository(Supabase)
	orderDetailRepo := repositories.NewOrderDetailRepository(Supabase)
	receiptRepo := repositories.NewReceiptRepository(Supabase)
	ipHC := httpclients.NewIPHttpClient()
	orderService := services.NewOrderService(orderRepo, orderDetailRepo,receiptRepo)
	ipSV := sv.NewIpService(ipHC)
	stripeSV := sv.NewStripeService(ipHC)
	gateways.NewHTTPGateway(app, ipSV, stripeSV,orderService)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
