package gateways

import (
	"encoding/json"
	"fmt"
	"go-fiber-template/src/domain/entities"
	"go-fiber-template/src/services"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// var latestPayment struct {
// 	sync.Mutex
// 	Amount   int64
// 	Currency string
// 	Success  bool
// }
var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}
func (h *HTTPGateway) InputPrice(ctx *fiber.Ctx) error {
	BodyPrice := new(entities.BodyPrice)
	if err := ctx.BodyParser(&BodyPrice); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "payload is incorrect"})
	}

	link, err := h.StripeService.StripeCreatePriceFromPayload(BodyPrice)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "Error to get link"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "get link success", Data: link})
}

func (h *HTTPGateway) TestWebhook(ctx *fiber.Ctx) error {
	// Parse Stripe webhook payload
	payload := ctx.Body()
	event := struct {
		Data struct {
			Object map[string]interface{} `json:"object"`
		} `json:"data"`
	}{}

	// Unmarshal the payload
	err := json.Unmarshal(payload, &event)
	if err != nil {
		// Log the error payload for better debugging
		fmt.Println("❌ Error parsing Stripe webhook payload:", err, string(payload))
		return ctx.Status(fiber.StatusUnauthorized).JSON(entities.ResponseModel{
			Message: "Unauthorized Webhook.",
		})
	}

	// Extract metadata from the event
	metadataInterface, exists := event.Data.Object["metadata"]
	if !exists {
		fmt.Println("❌ Error: missing metadata in webhook payload.")
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Missing metadata",
		})
	}

	metadata, ok := metadataInterface.(map[string]interface{})
	if !ok {
		fmt.Println("❌ Error: metadata is not a valid map.")
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Invalid metadata format",
		})
	}

	// Initialize the slice to hold order details
	var orderDetails []services.OrderDetailRequest

	// Extract the total item count (if provided)
	itemCount, countExists := metadata["item_count"].(string)
	if !countExists {
		fmt.Println("❌ Error: item_count is missing or invalid.")
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Missing item count",
		})
	}

	// Convert itemCount to integer
	itemCountInt, err := strconv.Atoi(itemCount)
	if err != nil {
		fmt.Println("❌ Error: invalid item count format:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Invalid item count format",
		})
	}

	// Loop through the metadata to extract item details dynamically
	for i := 0; i < itemCountInt; i++ {
		// Construct the keys dynamically for each item
		itemNameKey := fmt.Sprintf("item_name_%d", i)
		itemDescKey := fmt.Sprintf("item_desc_%d", i)
		itemPriceKey := fmt.Sprintf("item_price_%d", i)

		// Check if all necessary fields for the current item exist
		itemName, nameExists := metadata[itemNameKey].(string)
		itemDesc, descExists := metadata[itemDescKey].(string)
		itemPrice, priceExists := metadata[itemPriceKey].(string) // This should be a string

		// If any of the required fields are missing for the item, skip it
		if !nameExists || !descExists || !priceExists {
			fmt.Printf("❌ Missing required fields for item %d\n", i)
			continue
		}

		// Convert price to the correct type
		itemPriceInt, err := strconv.ParseInt(itemPrice, 10, 64) // Parse itemPrice as int64
		if err != nil {
			fmt.Println("❌ Error: invalid item price format:", err)
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
				Message: "Invalid item price format",
			})
		}

		// Construct the order detail
		orderDetail := services.OrderDetailRequest{
			Name:        itemName,
			UnitPrice:   itemPriceInt,                    // Set UnitPrice as int64
			Quantity:    1,                               // Since no quantity in metadata, assuming it's 1 by default
			TotalPrice:  itemPriceInt / 100,              // Set TotalPrice as itemPrice
			CreatedAt:   time.Now().Format(time.RFC3339), // Set the current timestamp
			Description: itemDesc,                        // Set description as string
		}

		// Append order detail to the list
		orderDetails = append(orderDetails, orderDetail)
	}

	// Check if no order details were found after processing
	if len(orderDetails) == 0 {
		fmt.Println("❌ Error: no valid order details found.")
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "No valid order details found",
		})
	}
	itemEmployee, employeeExists := metadata["employee"].(string)
	if !employeeExists {
		itemEmployee = ""
	}
	itemPayment, paymentExists := event.Data.Object["payment_method_types"]
	if !paymentExists {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Missing payment method types",
		})
	}
	paymentMethodTypes, ok := itemPayment.([]interface{})
	if !ok {
		fmt.Println("❌ Error: payment_method_types is not a valid array of interfaces.")
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Invalid payment_method_types format",
		})
	}

	// Convert []interface{} to []string
	var paymentMethodTypesStr []string
	for _, method := range paymentMethodTypes {
		if methodStr, ok := method.(string); ok {
			paymentMethodTypesStr = append(paymentMethodTypesStr, methodStr)
		} else {
			fmt.Println("❌ Error: payment_method_types contains a non-string value.")
			return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
				Message: "Invalid value in payment_method_types",
			})
		}
	}
	// Create the order request with the order details
	orderRequest := services.OrderRequest{
		OrderDetails:  orderDetails,
		PaymentMethod: paymentMethodTypesStr[0],
		Employee:      itemEmployee,
	}

	// Call OrderService to create the order with the details
	order, err := h.OrderService.CreateOrderWithDetails(orderRequest)
	if err != nil {
		fmt.Println("❌ Failed to create order:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: "Failed to create order",
		})
	}
	amountTotal, amountExists := event.Data.Object["amount_total"].(float64)
	if !amountExists {
		fmt.Println("❌ Error: missing amount in webhook payload.")
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "Missing amount",
		})
	}


	// Convert amount to proper format
	amount := int64(amountTotal) / 100
	message := fmt.Sprintf("ได้รับเงิน %d บาท", amount)

	// ✅ Broadcast payment update to WebSocket clients
	mutex.Lock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
	mutex.Unlock()

	fmt.Println("🔔 WebSocket Sent:", message)
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "Success",
		Data:    order,
	})
}

// func (h *HTTPGateway) GetPaymentStatus(ctx *fiber.Ctx) error {
// 	latestPayment.Lock()
// 	defer latestPayment.Unlock()

// 	// If no successful payment exists, return empty response
// 	if !latestPayment.Success {
// 		return ctx.JSON(fiber.Map{"success": false})
// 	}

// 	// Prepare response
// 	response := fiber.Map{
// 		"success":  true,
// 		"amount":   latestPayment.Amount,
// 		"currency": latestPayment.Currency,
// 	}

// 	// Reset payment status after sending response
// 	latestPayment.Success = false

// 	return ctx.JSON(response)
// }
func (h *HTTPGateway) WebSocketHandler(c *websocket.Conn) {
	mutex.Lock()
	clients[c] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, c)
		mutex.Unlock()
		c.Close()
	}()

	fmt.Println("✅ New WebSocket client connected")

	// ✅ Send a test message every 5 seconds
	for {
		time.Sleep(5 * time.Second)
	}
}