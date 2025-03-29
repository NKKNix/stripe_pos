package repositories

import (
	"fmt"

	supa "github.com/supabase-community/supabase-go"
)

type OrderDetail struct {
	ID          string `json:"id,omitempty"` // Ensure ID is optional
	OrderID     int    `json:"order_id"`     // Make sure this is populated
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unitprice"`
	TotalPrice  int64  `json:"totalprice"`
	Description string `json:"desc"`
}

type OrderDetailRepository struct {
	supabase *supa.Client
}

func NewOrderDetailRepository(supabase *supa.Client) *OrderDetailRepository {
	return &OrderDetailRepository{supabase: supabase}
}

// CreateOrderDetail inserts the order detail into the Supabase database
func (r *OrderDetailRepository) CreateOrderDetail(orderDetail *OrderDetail) (*OrderDetail, error) {

	// Ensure order_id is set before inserting
	if orderDetail.OrderID == 0 {
		return nil, fmt.Errorf("order_id must be set before creating order detail")
	}

	// Ensure created_at is managed by the database

	// Insert into Supabase (Correct return values)
	_, _, err := r.supabase.From("orderdetail").
		Insert([]OrderDetail{*orderDetail}, false, "", "representation", "minimal").
		Execute()

	// Enhanced error logging
	if err != nil {
		fmt.Println("❌ Error inserting order detail into database:", err)
		// Log the specific error from Supabase for further inspection
		fmt.Println("🔍 Supabase Error Details:", err.Error())
		return nil, err
	}

	return orderDetail, nil
}
