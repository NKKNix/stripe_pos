package entities

type Item struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       int64   `json:"price"` // Assuming price is in the smallest currency unit (e.g., satangs for THB)
    Image string `json:"image"`
}

type BodyPrice struct {
    Items    []Item   `json:"item"`
    Currency string   `json:"currency"`
    Method   []string `json:"method"`
    Employee string `json:"employee"`
}


