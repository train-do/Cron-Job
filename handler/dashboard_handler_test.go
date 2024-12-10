package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"project/handler"
	dashboardrepository "project/repository/dashboard_repository"
	"project/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func baseDahsboard() (handler.DashboardHandler, *dashboardrepository.DashboardRepoMock) {

	log := *zap.NewNop()

	mockService := &dashboardrepository.DashboardRepoMock{}
	service := service.Service{
		Dashboard: mockService,
	}

	return handler.NewDashboardHandler(&service, &log), mockService
}

func TestGetEarningDashboard(t *testing.T) {
	handler, mockRepo := baseDahsboard()
	t.Run("Successfully retrieve earning", func(t *testing.T) {

		r := gin.Default()
		r.GET("/dashboard/earning", handler.GetEarningDashboard)

		mockRepo.On("GetEarningDashboard").Once().Return(10000, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/earning", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertCalled(t, "GetEarningDashboard")

		var actualResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "successfully retrieved earning", actualResponse["message"])
		assert.True(t, actualResponse["status"].(bool))

		data := actualResponse["data"].(map[string]interface{})
		assert.Equal(t, float64(10000), data["total_earning"]) // JSON unmarshals numbers as float64
	})

	t.Run("Fail to retrieve earning - Service error", func(t *testing.T) {

		r := gin.Default()
		r.GET("/dashboard/earning", handler.GetEarningDashboard)

		mockRepo.On("GetEarningDashboard").Return(0, fmt.Errorf("no earnings")).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/earning", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockRepo.AssertCalled(t, "GetEarningDashboard")

		expectedResponse := `{"message":"There is no earning yet", "status":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}
