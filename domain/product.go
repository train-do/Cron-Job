package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"type:varchar(50);not null" json:"name" binding:"required,min=5"`
	SKUProduct  string          `gorm:"type:varchar(100);unique;not null" json:"sku_product" binding:"required"`
	Price       float64         `gorm:"not null" json:"price" binding:"required"`
	Description string          `gorm:"type:text;not null" json:"description" binding:"required"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`

	Image          []*Image          `gorm:"foreignKey:ProductID" json:"image"`
	ProductVariant []*ProductVariant `gorm:"foreignKey:ProductID" json:"product_variant"`
}

func SeedProducts() []Product {
	products := []Product{
		{
			Name:        "Product A",
			SKUProduct:  "SKU-001",
			Price:       100000,
			Description: "This is Product A description.",
		},
		{
			Name:        "Product B",
			SKUProduct:  "SKU-002",
			Price:       150000,
			Description: "This is Product B description.",
		},
		{
			Name:        "Product C",
			SKUProduct:  "SKU-003",
			Price:       200000,
			Description: "This is Product C description.",
		},
		{
			Name:        "Product D",
			SKUProduct:  "SKU-004",
			Price:       125000,
			Description: "This is Product D description.",
		},
		{
			Name:        "Product E",
			SKUProduct:  "SKU-005",
			Price:       175000,
			Description: "This is Product E description.",
		},
		{
			Name:        "Product F",
			SKUProduct:  "SKU-006",
			Price:       110000,
			Description: "This is Product F description.",
		},
		{
			Name:        "Product G",
			SKUProduct:  "SKU-007",
			Price:       130000,
			Description: "This is Product G description.",
		},
		{
			Name:        "Product H",
			SKUProduct:  "SKU-008",
			Price:       140000,
			Description: "This is Product H description.",
		},
		{
			Name:        "Product I",
			SKUProduct:  "SKU-009",
			Price:       160000,
			Description: "This is Product I description.",
		},
		{
			Name:        "Product J",
			SKUProduct:  "SKU-010",
			Price:       170000,
			Description: "This is Product J description.",
		},
		{
			Name:        "Product K",
			SKUProduct:  "SKU-011",
			Price:       180000,
			Description: "This is Product K description.",
		},
		{
			Name:        "Product L",
			SKUProduct:  "SKU-012",
			Price:       190000,
			Description: "This is Product L description.",
		},
		{
			Name:        "Product M",
			SKUProduct:  "SKU-013",
			Price:       210000,
			Description: "This is Product M description.",
		},
		{
			Name:        "Product N",
			SKUProduct:  "SKU-014",
			Price:       220000,
			Description: "This is Product N description.",
		},
		{
			Name:        "Product O",
			SKUProduct:  "SKU-015",
			Price:       230000,
			Description: "This is Product O description.",
		},
		{
			Name:        "Product P",
			SKUProduct:  "SKU-016",
			Price:       240000,
			Description: "This is Product P description.",
		},
		{
			Name:        "Product Q",
			SKUProduct:  "SKU-017",
			Price:       250000,
			Description: "This is Product Q description.",
		},
		{
			Name:        "Product R",
			SKUProduct:  "SKU-018",
			Price:       205000,
			Description: "This is Product R description.",
		},
		{
			Name:        "Product S",
			SKUProduct:  "SKU-019",
			Price:       105000,
			Description: "This is Product S description.",
		},
		{
			Name:        "Product T",
			SKUProduct:  "SKU-020",
			Price:       157500,
			Description: "This is Product T description.",
		},
		{
			Name:        "Product U",
			SKUProduct:  "SKU-021",
			Price:       215500,
			Description: "This is Product U description.",
		},
		{
			Name:        "Product V",
			SKUProduct:  "SKU-022",
			Price:       225750,
			Description: "This is Product V description.",
		},
		{
			Name:        "Product W",
			SKUProduct:  "SKU-023",
			Price:       150750,
			Description: "This is Product W description.",
		},
		{
			Name:        "Product X",
			SKUProduct:  "SKU-024",
			Price:       212500,
			Description: "This is Product X description.",
		},
		{
			Name:        "Product Y",
			SKUProduct:  "SKU-025",
			Price:       145750,
			Description: "This is Product Y description.",
		},
		{
			Name:        "Product Z",
			SKUProduct:  "SKU-026",
			Price:       190750,
			Description: "This is Product Z description.",
		},
	}

	return products
}

type BestSeller struct {
	ProductID int
	TotalSold int
}

type Revenue struct {
	Month   string `json:"month"`
	Revenue int    `json:"revenue"`
}
