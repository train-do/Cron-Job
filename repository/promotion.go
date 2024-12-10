package repository

import (
	"errors"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryPromotion interface {
	FindAll() ([]domain.Promotion, error)
	FindById(id uint) (domain.Promotion, error)
	Insert(stock *domain.Promotion) error
	Delete(stock *domain.Promotion) error
}

type repositoryPromotion struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewRepositoryPromotion(db *gorm.DB, log *zap.Logger) RepositoryPromotion {
	return &repositoryPromotion{db, log}
}

func (repo *repositoryPromotion) FindAll() ([]domain.Promotion, error) {
	promotions := []domain.Promotion{}
	if err := repo.db.Find(&promotions).Error; err != nil {
		// fmt.Println(err)
		return []domain.Promotion{}, errors.New(" Internal Server Error")
	}
	return promotions, nil
}
func (repo *repositoryPromotion) FindById(id uint) (domain.Promotion, error) {
	promotion := domain.Promotion{}
	if err := repo.db.Find(&promotion, "id=?", id).Error; err != nil {
		// fmt.Println(err)
		return domain.Promotion{}, errors.New(" Internal Server Error")
	}
	return promotion, nil
}
func (repo *repositoryPromotion) Insert(promotion *domain.Promotion) error {
	err := repo.db.Create(&promotion).Error
	if err != nil {
		// fmt.Println(err)
		return errors.New(" Bad Request")
	}
	return nil
}
func (repo *repositoryPromotion) Delete(promotion *domain.Promotion) error {
	if err := repo.db.Exec("DELETE FROM promotions WHERE id = ? RETURNING *", promotion.ID).Scan(&promotion).Error; err != nil {
		return errors.New(" Invalid ID")
	}
	return nil
}
