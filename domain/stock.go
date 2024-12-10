package domain

type Stock struct {
	ID               uint
	ProductVariantId int
	Description      string
	Qty              int
}
type ResponseStock struct {
	ProductName    string    `json:"product,omitempty"`
	ProductVariant SizeColor `json:"variant,omitempty"`
	Description    string    `json:"description,omitempty"`
	Qty            int       `json:"qty,omitempty"`
	CurrentStock   int       `json:"currentStock,omitempty"`
}
type SizeColor struct {
	Size  string `json:"size,omitempty"`
	Color string `json:"color,omitempty"`
}

func SeedStock() []Stock {

	stocks := []Stock{
		{
			ProductVariantId: 2,
			Description:      "Penambahan Manual",
			Qty:              3,
		},
		{
			ProductVariantId: 2,
			Description:      "Penambahan Manual",
			Qty:              10,
		},
	}

	return stocks
}
