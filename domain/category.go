package domain

import (
	"errors"
	"time"
)

type Category struct {
	ID        uint   `gorm:"primaryKey;autoincrement" json:"id"`
	Name      string `gorm:"type:varchar(50)" json:"name" binding:"required"`
	Icon      string `gorm:"type:varchar(200)" json:"image" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Icon == "" {
		return errors.New("icon is required")
	}
	return nil
}

func CategorySeeder() []Category {
	return []Category{
		{
			Name:      "Poultry",
			Icon:      "poultry-icon.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Livestock",
			Icon:      "livestock-icon.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Vegetables",
			Icon:      "vegetables-icon.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Fruits",
			Icon:      "fruits-icon.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
