package handler

import (
	"net/http"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardHandler interface {
	GetEarningDashboard(c *gin.Context)
	GetSummary(c *gin.Context)
	GetBestSeller(c *gin.Context)
	GetMonthlyRevenue(c *gin.Context)
}

type dashboardHandler struct {
	service *service.Service
	log     *zap.Logger
}

func NewDashboardHandler(service *service.Service, log *zap.Logger) DashboardHandler {
	return &dashboardHandler{service, log}
}

// GetEarningDashboard
// @Summary Retrieve total earning from dashboard
// @Description Get the total earning data from the dashboard service.
// @Tags Dashboard
// @Produce json
// @Success 200 {object} handler.Response(data=total_earning) "Success Response"
// @Failure 400 {object} handler.Response "Error Response"
// @Router /dashboard/earning [get]
func (dh *dashboardHandler) GetEarningDashboard(c *gin.Context) {
	totalEarning, err := dh.service.Dashboard.GetEarningDashboard()
	if err != nil {
		BadResponse(c, "There is no earning yet", http.StatusBadRequest)
		return
	}

	data := make(map[string]interface{})
	data["total_earning"] = totalEarning

	GoodResponseWithData(c, "successfully retrieved earning", http.StatusOK, data)
}

// @Summary Get summary of earnings
// @Description Retrieves the summary of earnings
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} handler.Response{data=domain.Summary} "Successfully retrieved summary"
// @Failure 400 {object} handler.Response "Error retrieving summary"
// @Router /dashboard/summary [get]
func (dh *dashboardHandler) GetSummary(c *gin.Context) {

	summary, err := dh.service.Dashboard.GetSummary()
	if err != nil {
		BadResponse(c, "There is no summary yet: "+err.Error(), http.StatusBadRequest)
		return
	}

	GoodResponseWithData(c, "successfully retrieved earning", http.StatusOK, summary)
}

// GetBestSeller
// @Summary Retrieve best seller products
// @Description Get the list of best seller products based on sales.
// @Tags Dashboard
// @Produce json
// @Success 200 {object} handler.Response{data=[]domain.BestSeller} "Success Response"
// @Failure 400 {object} handler.Response "Error Response"
// @Router /dashboard/bestSeller [get]
func (dh *dashboardHandler) GetBestSeller(c *gin.Context) {
	bestSellers, err := dh.service.Dashboard.GetBestSeller()
	if err != nil {
		BadResponse(c, "not found best seller: "+err.Error(), http.StatusBadRequest)
		return
	}

	GoodResponseWithData(c, "successfully retrieved", http.StatusOK, bestSellers)
}

// @Summary Get monthly revenue
// @Description Retrieves the monthly revenue
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} handler.Response{data=[]domain.Revenue} "Successfully retrieved monthly revenue"
// @Failure 400 {object} handler.Response "Error retrieving monthly revenue"
// @Router /dashboard/revenue [get]
func (dh *dashboardHandler) GetMonthlyRevenue(c *gin.Context) {

	revenue, err := dh.service.Dashboard.GetMonthlyRevenue()
	if err != nil {
		BadResponse(c, "Not Found revenue: "+err.Error(), http.StatusBadRequest)
		return
	}

	GoodResponseWithData(c, "successfully retrieved", http.StatusOK, revenue)
}
