package repositories

import (
	"fmt"

	supa "github.com/supabase-community/supabase-go"
)

type Receipt struct {
	ID            string `json:"id,omitempty"` // Ensure ID is optional
	OrderID       int    `json:"order_id"`     // Make sure this is populated
	PaymentMethod string `json:"payment_method"`
	Status        string    `json:"status"`
	TotalPrice    int64  `json:"total_amount"`
	Employee      string `json:"employee"`
}

type ReceiptRepository struct {
	supabase *supa.Client
}

func NewReceiptRepository(supabase *supa.Client) *ReceiptRepository {
	return &ReceiptRepository{supabase: supabase}
}

// CreateOrderDetail inserts the order detail into the Supabase database
func (r *ReceiptRepository) CreateReceipt(receipt *Receipt) (*Receipt, error) {

	// Ensure order_id is set before inserting
	if receipt.OrderID == 0 {
		return nil, fmt.Errorf("order_id must be set before creating order detail")
	}

	// Ensure created_at is managed by the database

	// Insert into Supabase (Correct return values)
	_, _, err := r.supabase.From("receipt").
		Insert([]Receipt{*receipt}, false, "", "representation", "minimal").
		Execute()

	// Enhanced error logging
	if err != nil {
		fmt.Println("❌ Error inserting receipt into database:", err)
		// Log the specific error from Supabase for further inspection
		fmt.Println("🔍 Supabase Error Details:", err.Error())
		return nil, err
	}

	return receipt, nil
}
