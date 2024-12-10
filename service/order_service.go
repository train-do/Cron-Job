package service

import (
	"project/domain"
	"project/repository"
)

type OrderService interface {
	All(page, limit uint) (int, int, []domain.OrderTotal, error)
	Update(orderId uint, confirmation domain.OrderConfirmation) error
	Get(orderId uint) (domain.OrderTotal, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) All(page, limit uint) (int, int, []domain.OrderTotal, error) {
	return s.repo.All(page, limit)
}

func (s *orderService) Update(orderId uint, confirmation domain.OrderConfirmation) error {
	return s.repo.Update(orderId, confirmation)
}

func (s *orderService) Get(orderId uint) (domain.OrderTotal, error) {
	return s.repo.Get(orderId)
}
