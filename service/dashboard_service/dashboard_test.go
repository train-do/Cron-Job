package dashboardservice_test

import (
	"fmt"
	"project/domain"
	"project/repository"
	dashboardrepository "project/repository/dashboard_repository"
	dashboardservice "project/service/dashboard_service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func base() (dashboardservice.DashboardService, *dashboardrepository.DashboardRepoMock) {
	log := *zap.NewNop()
	mockRepo := &dashboardrepository.DashboardRepoMock{}
	repo := repository.Repository{
		Dashboard: mockRepo,
	}
	service := dashboardservice.NewDashboardService(&repo, &log)

	return service, mockRepo
}

func TestGetEarningDashboard(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully get earnings from product", func(t *testing.T) {
		expectedEarnings := 1000

		mockRepo.On("GetEarningDashboard").
			Return(expectedEarnings, nil).
			Once()

		totalEarnings, err := service.GetEarningDashboard()

		assert.NoError(t, err)
		assert.Equal(t, expectedEarnings, totalEarnings)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get earnings from product - Repository error", func(t *testing.T) {

		mockRepo.On("GetEarningDashboard").
			Return(0, fmt.Errorf("database error")).
			Once()

		totalEarnings, err := service.GetEarningDashboard()

		assert.Error(t, err)
		assert.Equal(t, 0, totalEarnings)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetSummary(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully get summary", func(t *testing.T) {
		expectedSummary := &domain.Summary{
			Sales:  1000,
			Orders: 10,
			Items:  50,
			Users:  5,
		}

		mockRepo.On("GetSummary").
			Return(expectedSummary, nil).
			Once()

		summary, err := service.GetSummary()

		assert.NoError(t, err)
		assert.Equal(t, expectedSummary, summary)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get summary - Repository error", func(t *testing.T) {
		mockRepo.On("GetSummary").
			Return(nil, fmt.Errorf("database error")).
			Once()

		summary, err := service.GetSummary()

		assert.Error(t, err)
		assert.Nil(t, summary)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetBestSeller(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully get best sellers", func(t *testing.T) {
		expectedBestSellers := []*domain.BestSeller{
			{ProductID: 1, TotalSold: 100},
			{ProductID: 2, TotalSold: 50},
		}

		mockRepo.On("GetBestSeller").
			Return(expectedBestSellers, nil).
			Once()

		bestSellers, err := service.GetBestSeller()

		assert.NoError(t, err)
		assert.Equal(t, expectedBestSellers, bestSellers)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get best sellers - Repository error", func(t *testing.T) {
		mockRepo.On("GetBestSeller").
			Return(nil, fmt.Errorf("database error")).
			Once()

		bestSellers, err := service.GetBestSeller()

		assert.Error(t, err)
		assert.Nil(t, bestSellers)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetMonthlyRevenue(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully get monthly revenue", func(t *testing.T) {
		expectedRevenue := []*domain.Revenue{
			{Month: "January", Revenue: 1000},
			{Month: "February", Revenue: 1500},
		}

		mockRepo.On("GetMonthlyRevenue").
			Return(expectedRevenue, nil).
			Once()

		revenue, err := service.GetMonthlyRevenue()

		assert.NoError(t, err)
		assert.Equal(t, expectedRevenue, revenue)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get monthly revenue - Repository error", func(t *testing.T) {
		mockRepo.On("GetMonthlyRevenue").
			Return(nil, fmt.Errorf("database error")).
			Once()

		revenue, err := service.GetMonthlyRevenue()

		assert.Error(t, err)
		assert.Nil(t, revenue)

		mockRepo.AssertExpectations(t)
	})
}
