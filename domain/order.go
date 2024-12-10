package domain

import (
	"errors"
	"time"
)

type Status string

const (
	Created   Status = "created"
	Processed Status = "processed"
	Canceled  Status = "canceled"
	Completed Status = "completed"
)

type Order struct {
	ID             uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID     uint        `json:"-"`
	Customer       Customer    `json:"customer"`
	PaymentMethod  string      `json:"payment_method"`
	TrackingNumber string      `json:"tracking_number"`
	Status         Status      `gorm:"type:orderstatus" json:"status"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (order *Order) Confirm(confirmation OrderConfirmation) error {
	currentStatus := map[Status]func(OrderConfirmation) error{
		Created:   order.Process,
		Processed: order.Ship,
	}

	if nextStatus, ok := currentStatus[order.Status]; ok {
		return nextStatus(confirmation)
	}

	return errors.New("invalid status")
}

func (order *Order) Process(confirmation OrderConfirmation) error {
	order.Status = Canceled
	if confirmation.Accept {
		order.Status = Processed
	}
	return nil
}

func (order *Order) Ship(confirmation OrderConfirmation) error {
	order.Status = Canceled
	if confirmation.Accept {
		order.Status = Completed
		return order.setTrackingNumber(*confirmation.TrackingNumber)
	}
	return nil
}

func (order *Order) setTrackingNumber(trackingNumber string) error {
	if trackingNumber == "" {
		return errors.New("invalid tracking number")
	}
	order.TrackingNumber = trackingNumber
	return nil
}
