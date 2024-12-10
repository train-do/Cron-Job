package domain

type status string

const (
	Active   status = "Active"
	Inactive status = "Inactive"
)

type Type string

const (
	Voucher  Type = "Voucher Code"
	Discount Type = "Direct Discount"
)

type Promotion struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	StartDate   string `gorm:"type:date"`
	EndDate     string `gorm:"type:date"`
	Type        Type
	Status      status
	IsPublish   bool
	VoucherCode string
	Limit       int
}

func SeedPromotions() []Promotion {
	promotions := []Promotion{
		{
			Name:        "Promo Akhir Tahun",
			Description: "Potongan 20'%' dengan pembelian di atas 100rb",
			StartDate:   "2024-12-15",
			EndDate:     "2024-12-22",
			Type:        "Direct Discount",
			Status:      "Active",
			IsPublish:   false,
			Limit:       20,
		},
		{
			Name:        "Cuci Gudang",
			Description: "Potongan 30'%' dengan pembelian di atas 100rb",
			StartDate:   "2024-12-15",
			EndDate:     "2024-12-22",
			Type:        "Voucher Code",
			Status:      "Active",
			IsPublish:   true,
			VoucherCode: "CGDG",
			Limit:       20,
		},
		{
			Name:        "Spesial Kemerdekaan",
			Description: "Potongan 10'%' dengan pembelian di atas 100rb",
			StartDate:   "2024-12-15",
			EndDate:     "2024-12-22",
			Type:        "Direct Discount",
			Status:      "Inactive",
			IsPublish:   false,
			Limit:       20,
		},
		{
			Name:        "Hari Kartini",
			Description: "Potongan 15'%' dengan pembelian di atas 100rb",
			StartDate:   "2024-12-15",
			EndDate:     "2024-12-22",
			Type:        "Direct Discount",
			Status:      "Inactive",
			IsPublish:   false,
			Limit:       20,
		},
	}

	return promotions
}
