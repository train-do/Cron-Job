package domain

type OrderItemSubtotal struct {
	OrderID     uint    `json:"-"`
	ProductName string  `json:"product_name"`
	Quantity    uint    `json:"quantity"`
	UnitPrice   float32 `json:"unit_price"`
	Subtotal    float32 `json:"subtotal"`
}
