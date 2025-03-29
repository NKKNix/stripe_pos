package repositories

import (
	"encoding/json"
	"fmt"

	supa "github.com/supabase-community/supabase-go"
)

type Order struct {
	ID          string  `json:"id,omitempty"` // Ensure ID is optional
	OrderDate   string  `json:"created_at"`
	TotalAmount int64 `json:"TotalAmount"`
	Status      string  `json:"Status"`
}

type OrderRepository struct {
	supabase *supa.Client
}

func NewOrderRepository(supabase *supa.Client) *OrderRepository {
	return &OrderRepository{supabase: supabase}
}

// Create Order
func (r *OrderRepository) CreateOrder(order *Order) (int, error) {


	// Insert into Supabase
	data, _, err := r.supabase.From("order").
		Insert([]Order{*order}, false, "", "representation", "minimal").Execute()

	if err != nil {
		fmt.Println("❌ Error inserting order into database:", err)
		return 0, err
	}

	// If data is returned as []byte, unmarshal it
	var response []map[string]interface{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("❌ Error unmarshalling response:", err)
		return 0, err
	}
	createID := response[0]["id"]
    fmt.Println("✅ Order ID Successfully:", createID)

    // Debugging the type of createID
    fmt.Printf("Type of createID: %T\n", createID)

    // Assert createID to the appropriate type based on your actual data
    switch v := createID.(type) {
    case int:
        // If it's an integer
        return int(v), nil
    case float64:
        // If it's a float (common for JSON numbers)
        return int(v), nil
    default:
        fmt.Println("Error: value is not an int or float")
        return 0, fmt.Errorf("expected int or float, got %T", v)
    }
}
