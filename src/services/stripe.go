package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fiber-template/src/domain/entities"
	httpclients "go-fiber-template/src/infrastructure/httpclients"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

type StripeService struct {
	IpHttpClient httpclients.IpInterface
}

// NewStripeService creates a new StripeService instance
func NewStripeService(ipHC httpclients.IpInterface) IStripeService {
	return &StripeService{
		IpHttpClient: ipHC,
	}
}

type IStripeService interface {
	// StripeCreatePrice(userID string, BodyPrice *entities.BodyPrice) (string, error)
	StripeCreatePriceFromPayload(payload *entities.BodyPrice) (string, error)
	ConfirmMoney(totalAmount int) ([]byte, error)
}

func (sv StripeService) StripeCreatePriceFromPayload(payload *entities.BodyPrice) (string, error) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	url := os.Getenv("STRIPE_REDIRECT")

	var lineItems []*stripe.CheckoutSessionLineItemParams

	// Iterate through items array and build line items
	for _, item := range payload.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(payload.Currency),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(item.Name),
					Description: stripe.String(item.Description),
					Images: []*string{
						stripe.String(item.Image),
					},
				},
				UnitAmount: stripe.Int64(item.Price),
			},
			Quantity: stripe.Int64(1),
		})

	}
	metadata := map[string]string{
		"item_count": fmt.Sprintf("%d", len(payload.Items)),
	}
	metadata["employee"] = payload.Employee

	// Add item details to metadata (optional, limited by Stripe's metadata constraints)
	for i, item := range payload.Items {
		metadata[fmt.Sprintf("item_name_%d", i)] = item.Name
		metadata[fmt.Sprintf("item_price_%d", i)] = fmt.Sprintf("%d", item.Price)
		metadata[fmt.Sprintf("item_desc_%d", i)] = item.Description
	}

	// Stripe checkout session parameters
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes:  stripe.StringSlice(payload.Method),
		LineItems:           lineItems, // Attach dynamic line items
		Mode:                stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:          stripe.String(url),
		CancelURL:           stripe.String(os.Getenv("FRONT_REDIRECT_URL_STRIPE")),
		AllowPromotionCodes: stripe.Bool(true),
		ExpiresAt:           stripe.Int64(time.Now().Add(60 * time.Minute).Unix()),
		Metadata:            metadata, // Attach metadata here
	}

	session, err := session.New(params)
	if err != nil {
		return "", err
	}

	// Return structured response
	return session.URL, nil
}

func (sv StripeService) ConfirmMoney(totalAmount int) ([]byte, error) {
	url := "https://api-voice.botnoi.ai/openapi/v1/generate_audio"
	method := "POST"
	fmt.Println("this is confirm price: ", totalAmount)
	text := fmt.Sprintf("จำนวนเงิน %d บาท", totalAmount/100)
	var payload = map[string]interface{}{"text": text, "speaker": "1", "volume": 1, "speed": 1, "type_media": "m4a", "save_file": "false", "language": "th"}
	jsonValue, _ := json.Marshal(payload)

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))

	req.Header.Add("Botnoi-Token", "UXpKT1FrUEZKY1FuU2lBUmU0bVI4czN6MkV6MTU2MTg5NA==")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	return body, nil
}
