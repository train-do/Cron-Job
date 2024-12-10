package domain

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        int             `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID int             `gorm:"not null" json:"product_id"`
	URLPath   string          `gorm:"type:varchar(150)" json:"url_path"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

func SeedImages() []Image {
	images := []Image{
		{
			ProductID: 1,
			URLPath:   "https://example.com/images/product/aneka_buah1.jpg",
		},
		{
			ProductID: 1,
			URLPath:   "https://example.com/images/product/aneka_buah2.jpg",
		},
		{
			ProductID: 1,
			URLPath:   "https://example.com/images/product/aneka_buah3.jpg",
		},
		{
			ProductID: 2,
			URLPath:   "https://example.com/images/product/aneka_mangga1.jpg",
		},
		{
			ProductID: 2,
			URLPath:   "https://example.com/images/product/aneka_mangga2.jpg",
		},
		{
			ProductID: 2,
			URLPath:   "https://example.com/images/product/aneka_mangga3.jpg",
		},
		{
			ProductID: 3,
			URLPath:   "https://example.com/images/product/aneka_apel1.jpg",
		},
		{
			ProductID: 3,
			URLPath:   "https://example.com/images/product/aneka_apel2.jpg",
		},
		{
			ProductID: 3,
			URLPath:   "https://example.com/images/product/aneka_apel3.jpg",
		},
	}

	return images
}
