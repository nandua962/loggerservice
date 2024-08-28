package entities

// PaymentGatewayName represents a simple struct to hold the name of a payment gateway.
type PaymentGatewayName struct {
	Name string `json:"name"`
}

// PaymentGatewayName represents a simple struct to hold the name of a payment gateway.
type PaymentGateway struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
