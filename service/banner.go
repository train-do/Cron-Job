package service

import (
	"project/domain"
	"project/repository"
)

type ServiceBanner interface {
	GetAll() ([]domain.Banner, error)
	GetById(id uint) (domain.Banner, error)
	Create(banner *domain.Banner) error
	Edit(banner *domain.Banner) error
	Delete(banner *domain.Banner) error
}

type serviceBanner struct {
	repo repository.RepositoryBanner
}

func NewServiceBanner(repo repository.RepositoryBanner) ServiceBanner {
	return &serviceBanner{repo: repo}
}

func (s *serviceBanner) GetAll() ([]domain.Banner, error) {
	return s.repo.FindAll()
}
func (s *serviceBanner) GetById(id uint) (domain.Banner, error) {
	return s.repo.FindById(id)
}
func (s *serviceBanner) Create(banner *domain.Banner) error {
	return s.repo.Insert(banner)
}
func (s *serviceBanner) Edit(banner *domain.Banner) error {
	return s.repo.Update(banner)
}
func (s *serviceBanner) Delete(banner *domain.Banner) error {
	return s.repo.Delete(banner)
}
