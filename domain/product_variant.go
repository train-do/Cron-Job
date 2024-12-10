package domain

import (
	"errors"
	"log"
	"time"

	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type ProductVariant struct {
	ID        int             `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID int             `gorm:"not null" json:"product_id"`
	Size      string          `gorm:"type:varchar(50)" json:"size"`
	Color     string          `gorm:"type:varchar(50)" json:"color"`
	Stock     int             `gorm:"default:0;check:stock>=0" json:"stock"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

func (variant *ProductVariant) DeductStock(quantity uint) error {
	qty := int(quantity)
	log.Println("before", variant.Stock)
	if variant.Stock >= qty {
		variant.Stock -= qty
		log.Println("after", variant.Stock)
		return nil
	}

	return errors.New("not enough stock")
}

func SeedProductVariants() []ProductVariant {
	sizes := []string{"S", "M", "L", "XL"}
	colors := []string{"Red", "Blue", "Green"}
	var variants []ProductVariant

	for productID := 1; productID <= 26; productID++ {
		for _, size := range sizes {
			for _, color := range colors {
				variants = append(variants, ProductVariant{
					ProductID: productID,
					Size:      size,
					Color:     color,
					Stock:     rand.Intn(50) + 1,
				})
			}
		}
	}

	return variants
}
