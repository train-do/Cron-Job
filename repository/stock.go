package repository

import (
	"errors"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryStock interface {
	FindById(id int) (domain.ResponseStock, error)
	Insert(stock *domain.Stock) error
	Delete(stock *domain.Stock) error
}

type repositoryStock struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewRepositoryStock(db *gorm.DB, log *zap.Logger) RepositoryStock {
	return &repositoryStock{db, log}
}

func (repo *repositoryStock) FindById(id int) (domain.ResponseStock, error) {
	var productVariant domain.ProductVariant
	if err := repo.db.Find(&productVariant, "id=?", id).Error; err != nil {
		// fmt.Println(err)
		return domain.ResponseStock{}, errors.New(" Invalid ID")
	}
	var product domain.Product
	if err := repo.db.Find(&product, "id=?", productVariant.ProductID).Error; err != nil {
		// fmt.Println(err)
		return domain.ResponseStock{}, errors.New(" Invalid ID")
	}
	result := domain.ResponseStock{
		ProductName:    product.Name,
		ProductVariant: domain.SizeColor{Size: productVariant.Size, Color: productVariant.Color},
		CurrentStock:   productVariant.Stock,
	}
	return result, nil
}
func (repo *repositoryStock) Insert(stock *domain.Stock) error {
	err := repo.db.Create(&stock).Error
	if err != nil {
		// fmt.Println(err)
		return errors.New(" Bad Request")
	}
	return nil
}
func (repo *repositoryStock) Delete(stock *domain.Stock) error {
	if err := repo.db.First(&stock, stock.ID).Error; err != nil {
		return errors.New(" Invalid ID")
	}
	repo.db.Delete(&stock)
	return nil
}
