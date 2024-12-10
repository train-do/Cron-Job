package service

import (
	"project/domain"
	"project/repository"
)

type ServicePromotion interface {
	GetAll() ([]domain.Promotion, error)
	GetById(id uint) (domain.Promotion, error)
	Create(promotion *domain.Promotion) error
	Delete(promotion *domain.Promotion) error
}

type servicePromotion struct {
	repo repository.RepositoryPromotion
}

func NewServicePromotion(repo repository.RepositoryPromotion) ServicePromotion {
	return &servicePromotion{repo: repo}
}

func (s *servicePromotion) GetAll() ([]domain.Promotion, error) {
	return s.repo.FindAll()
}
func (s *servicePromotion) GetById(id uint) (domain.Promotion, error) {
	return s.repo.FindById(id)
}
func (s *servicePromotion) Create(promotion *domain.Promotion) error {
	return s.repo.Insert(promotion)
}
func (s *servicePromotion) Delete(promotion *domain.Promotion) error {
	return s.repo.Delete(promotion)
}
