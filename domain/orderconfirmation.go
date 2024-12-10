package domain

type OrderConfirmation struct {
	Accept         bool    `json:"accept"`
	TrackingNumber *string `json:"tracking_number"`
}
