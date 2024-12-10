package repository

import (
	"errors"
	"gorm.io/gorm"
	"project/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo UserRepository) Create(user *domain.User) error {
	return repo.db.Create(&user).Error
}

func (repo UserRepository) All(user domain.User) ([]domain.User, error) {
	var users []domain.User
	result := repo.db.Where(user).Find(&users)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return users, nil
}
