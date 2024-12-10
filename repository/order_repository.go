package repository

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"project/domain"
	"project/helper"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo OrderRepository) Update(orderId uint, confirmation domain.OrderConfirmation) error {
	var order = domain.Order{ID: orderId}

	repo.db.Preload("Items.Variant").First(&order)
	if err := order.Confirm(confirmation); err != nil {
		return err
	}

	if err := repo.shouldUpdateStock(&order); err != nil {
		return err
	}

	return repo.db.Save(&order).Error
}

func (repo OrderRepository) shouldUpdateStock(order *domain.Order) error {
	if order.Status == domain.Processed {
		for _, item := range order.Items {
			if err := item.Variant.DeductStock(item.Quantity); err != nil {
				return err
			}
			repo.db.Save(&item.Variant)
		}
	}
	return nil
}

func (repo OrderRepository) All(page, limit uint) (int, int, []domain.OrderTotal, error) {
	var count int64
	repo.db.Model(&domain.OrderTotal{}).Count(&count)
	pages := int(math.Ceil(float64(count) / float64(limit)))

	var orders []domain.OrderTotal
	result := repo.db.Scopes(helper.Paginate(page, limit)).Find(&orders)
	if result.RowsAffected == 0 {
		return 0, 0, nil, errors.New("order not found")
	}
	return int(count), pages, orders, nil
}

func (repo OrderRepository) Get(orderId uint) (domain.OrderTotal, error) {
	var order domain.OrderTotal
	result := repo.db.Preload("Items").First(&order, orderId)
	return order, result.Error
}
