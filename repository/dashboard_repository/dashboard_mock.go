package dashboardrepository

import (
	"project/domain"

	"github.com/stretchr/testify/mock"
)

type DashboardRepoMock struct {
	mock.Mock
}

func (pr *DashboardRepoMock) GetEarningDashboard() (int, error) {
	args := pr.Called()
	return args.Int(0), args.Error(1)
}

func (pr *DashboardRepoMock) GetSummary() (*domain.Summary, error) {
	args := pr.Called()
	if summary, ok := args.Get(0).(*domain.Summary); ok {
		return summary, args.Error(1)
	}
	return nil, args.Error(1)
}

func (pr *DashboardRepoMock) GetBestSeller() ([]*domain.BestSeller, error) {
	args := pr.Called()
	if bestSellers, ok := args.Get(0).([]*domain.BestSeller); ok {
		return bestSellers, args.Error(1)
	}
	return nil, args.Error(1)
}

func (pr *DashboardRepoMock) GetMonthlyRevenue() ([]*domain.Revenue, error) {
	args := pr.Called()
	if revenue, ok := args.Get(0).([]*domain.Revenue); ok {
		return revenue, args.Error(1)
	}
	return nil, args.Error(1)
}
