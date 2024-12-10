package service

import (
	"project/domain"
	"project/repository"
)

type UserService interface {
	All(user domain.User) ([]domain.User, error)
	Register(user *domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) All(user domain.User) ([]domain.User, error) {
	return s.repo.All(user)
}

func (s *userService) Register(user *domain.User) error {
	return s.repo.Create(user)
}
