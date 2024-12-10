package dashboardservice

import (
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type DashboardService interface {
	GetEarningDashboard() (int, error)
	GetSummary() (*domain.Summary, error)
	GetBestSeller() ([]*domain.BestSeller, error)
	GetMonthlyRevenue() ([]*domain.Revenue, error)
}

type dashboardService struct {
	repo *repository.Repository
	log  *zap.Logger
}

func NewDashboardService(repo *repository.Repository, log *zap.Logger) DashboardService {
	return &dashboardService{repo, log}
}

func (ds *dashboardService) GetEarningDashboard() (int, error) {

	totalEarning, err := ds.repo.Dashboard.GetEarningDashboard()
	if err != nil {
		return 0, err
	}

	return totalEarning, nil
}

func (ds *dashboardService) GetSummary() (*domain.Summary, error) {

	summary, err := ds.repo.Dashboard.GetSummary()
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (ds *dashboardService) GetBestSeller() ([]*domain.BestSeller, error) {

	bestSellers, err := ds.repo.Dashboard.GetBestSeller()
	if err != nil {
		return nil, err
	}

	return bestSellers, nil
}

func (ds *dashboardService) GetMonthlyRevenue() ([]*domain.Revenue, error) {

	revenue, err := ds.repo.Dashboard.GetMonthlyRevenue()
	if err != nil {
		return nil, err
	}

	return revenue, nil
}
