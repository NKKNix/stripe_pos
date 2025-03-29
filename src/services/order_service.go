package services

import (
	"fmt"
	"time"

	"go-fiber-template/src/domain/repositories"
)

// Define the interface for OrderService
type IOrderService interface {
	CreateOrderWithDetails(orderReq OrderRequest) (*repositories.Order, error)
}

// Implement the OrderService struct
type OrderService struct {
	orderRepo       *repositories.OrderRepository
	orderDetailRepo *repositories.OrderDetailRepository
	receiptRepo     *repositories.ReceiptRepository
}

// Ensure OrderService implements IOrderService

// Constructor for OrderService
func NewOrderService(orderRepo *repositories.OrderRepository, orderDetailRepo *repositories.OrderDetailRepository, receiptRepo *repositories.ReceiptRepository) IOrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		orderDetailRepo: orderDetailRepo,
		receiptRepo:     receiptRepo,
	}
}

// OrderRequest represents an order with multiple order details
type OrderRequest struct {
	OrderDetails []OrderDetailRequest `json:"order_details"`
	PaymentMethod string `json:"payment_method"`
	Employee      string `json:"employee"`
}

// OrderDetailRequest represents a single order detail
type OrderDetailRequest struct {
	ID            string `json:"id"`
	OrderID       int    `json:"order_id"`
	Name          string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	UnitPrice     int64  `json:"unit_price"`  // Changed to int64
	TotalPrice    int64  `json:"total_price"` // Changed to int64
	CreatedAt     string `json:"created_at"`  // Should be a formatted timestamp string
	Description   string `json:"description"` // Ensure this field is nullable in database if not used

}

// CreateOrderWithDetails - Creates an order and its details
func (s *OrderService) CreateOrderWithDetails(orderReq OrderRequest) (*repositories.Order, error) {
	fmt.Println("🟡 Received Order Request:", orderReq)

	// Calculate Total Amount
	var totalAmount int64
	for _, detail := range orderReq.OrderDetails {
		totalAmount += int64(detail.Quantity) * detail.UnitPrice
	}

	fmt.Println("🟡 Calculated Total Amount:", totalAmount)

	// Create Order
	order := &repositories.Order{
		OrderDate:   time.Now().Format("2006-01-02 15:04:05"),
		TotalAmount: totalAmount / 100,
		Status:      "Processing",
	}

	createdOrder, err := s.orderRepo.CreateOrder(order)
	fmt.Println("🟡 Created Order:", createdOrder)
	if err != nil {
		fmt.Println("❌ Error inserting order:", err)
		return nil, err
	}
	fmt.Println("✅ Order Inserted Successfully:", createdOrder)

	// Insert Order Details using OrderDetailRepo

	for _, detail := range orderReq.OrderDetails {
		orderDetail := &repositories.OrderDetail{
			OrderID:     createdOrder,                                    // Linking the order with the order detail
			Name:        detail.Name,                                     // Assigning product name
			UnitPrice:   detail.UnitPrice / 100,                          // Assigning unit price
			Quantity:    detail.Quantity,                                 // Assigning quantity
			TotalPrice:  int64(detail.Quantity) * detail.UnitPrice / 100, // Calculating total price
			Description: detail.Description,                              // Assigning the current timestamp
		}

		// Insert order details through the OrderDetailRepo
		_, err := s.orderDetailRepo.CreateOrderDetail(orderDetail)
		if err != nil {
			fmt.Println("❌ Error inserting order detail:", err)
			return nil, err
		}

	}
	receipt := &repositories.Receipt{
		OrderID:       createdOrder,
		PaymentMethod: orderReq.PaymentMethod, // Assuming PaymentMethod comes from the request
		Status:        "Processing",                      // You can modify this depending on your status codes
		TotalPrice:    totalAmount / 100,      // Use the total amount from order details
		Employee:      orderReq.Employee,      // Assuming Employee comes from the request
	}

	createdReceipt, err := s.receiptRepo.CreateReceipt(receipt)
	if err != nil {
		fmt.Println("❌ Error inserting receipt:", err)
		return nil, err
	}
	fmt.Println("✅ Receipt Created Successfully:", createdReceipt)
	return order, nil
}
