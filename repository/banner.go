package repository

import (
	"errors"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryBanner struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewRepositoryBanner(db *gorm.DB, log *zap.Logger) *RepositoryBanner {
	return &RepositoryBanner{db: db, log: log}
}

func (repo *RepositoryBanner) FindAll() ([]domain.Banner, error) {
	var banners []domain.Banner
	if err := repo.db.Find(&banners).Error; err != nil {
		// fmt.Println(err)
		return []domain.Banner{}, errors.New(" Internal Server Error")
	}
	return banners, nil
}
func (repo *RepositoryBanner) FindById(id uint) (domain.Banner, error) {
	var banner domain.Banner
	if err := repo.db.Where("id=?", id).Find(&banner).Error; err != nil {
		// fmt.Println(err)
		return domain.Banner{}, errors.New(" ID Banner Not Found")
	}
	return banner, nil
}
func (repo *RepositoryBanner) Insert(banner *domain.Banner) error {
	err := repo.db.Create(&banner).Error
	if err != nil {
		// fmt.Println(err)
		return errors.New(" Something Wrong")
	}
	return nil
}
func (repo *RepositoryBanner) Update(banner *domain.Banner) error {
	err := repo.db.Save(&banner).Error
	if err != nil {
		// fmt.Println(err)
		return errors.New(" Something Wrong")
	}
	return nil
}
func (repo *RepositoryBanner) Delete(banner *domain.Banner) error {
	if err := repo.db.First(&banner, banner.ID).Error; err != nil {
		return errors.New(" ID Banner Not Found")
	}
	err := repo.db.Delete(&banner).Error
	if err != nil {
		// fmt.Println(err)
		return errors.New(" Something Wrong")
	}
	return nil
}
