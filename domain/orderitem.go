package domain

type OrderItem struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"-"`
	OrderID   uint           `gorm:"not null" json:"-"`
	Order     Order          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	VariantID uint           `json:"-"`
	Variant   ProductVariant `json:"variant"`
	Quantity  uint           `json:"quantity"`
	UnitPrice float64        `gorm:"type:float" json:"unit_price"`
}
