package service

import (
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type ServiceStock interface {
	GetDetails(id int) (domain.ResponseStock, error)
	Edit(id int, newStock int) (domain.Stock, error)
	Delete(stock *domain.Stock) error
}

type serviceStock struct {
	repo repository.RepositoryStock
	log  *zap.Logger
}

func NewServiceStock(repo repository.RepositoryStock, log *zap.Logger) ServiceStock {
	return &serviceStock{repo, log}
}

func (service *serviceStock) GetDetails(id int) (domain.ResponseStock, error) {
	return service.repo.FindById(id)
}
func (service *serviceStock) Edit(id int, newStock int) (domain.Stock, error) {
	result, err := service.repo.FindById(id)
	if err != nil {
		return domain.Stock{}, err
	}
	var stock domain.Stock
	stock.ProductVariantId = id
	if result.CurrentStock > newStock {
		stock.Qty = result.CurrentStock - newStock
		stock.Description = "Pengurangan Manual"
	} else if result.CurrentStock < newStock {
		stock.Qty = newStock - result.CurrentStock
		stock.Description = "Penambahan Manual"
	}
	if err := service.repo.Insert(&stock); err != nil {
		return domain.Stock{}, err
	}
	return stock, nil
}
func (service *serviceStock) Delete(stock *domain.Stock) error {
	return service.repo.Delete(stock)
}
