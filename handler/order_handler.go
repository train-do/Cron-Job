package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"
)

type OrderController struct {
	service service.OrderService
	logger  *zap.Logger
}

func NewOrderController(service service.OrderService, logger *zap.Logger) *OrderController {
	return &OrderController{service: service, logger: logger}
}

// Order endpoint
// @Summary Customer orders
// @Description Get customer orders
// @Tags Order
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "orders retrieved"
// @Failure 404 {object} handler.Response "no data found"
// @Failure 500 {object} handler.Response "server error"
// @Router  /orders [get]
func (ctrl *OrderController) All(c *gin.Context) {
	page, _ := helper.Uint(c.Query("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := helper.Uint(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	total, pages, orders, err := ctrl.service.All(page, limit)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithPage(c, "orders retrieved", http.StatusOK, total, pages, int(page), int(limit), orders)
}

// Order endpoint
// @Summary Customer order
// @Description Update customer order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param id path uint true "Order ID"
// @Success 200 {object} handler.Response "order updated"
// @Failure 422 {object} handler.Response "invalid input"
// @Failure 404 {object} handler.Response "no data found"
// @Failure 500 {object} handler.Response "server error"
// @Router  /orders/:id [put]
func (ctrl *OrderController) Update(c *gin.Context) {
	orderId, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "invalid id", http.StatusUnprocessableEntity)
		return
	}

	var confirmation domain.OrderConfirmation
	if err := c.ShouldBindJSON(&confirmation); err != nil {
		BadResponse(c, "invalid input", http.StatusBadRequest)
		return
	}

	if err = ctrl.service.Update(orderId, confirmation); err != nil {
		BadResponse(c, "server error", http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "order status updated", http.StatusOK, nil)
}

// Order endpoint
// @Summary Customer order
// @Description Get customer order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param id path uint true "Order ID"
// @Success 200 {object} handler.Response "order retrived"
// @Failure 422 {object} handler.Response "invalid input"
// @Failure 404 {object} handler.Response "no data found"
// @Failure 500 {object} handler.Response "server error"
// @Router  /orders/:id [get]
func (ctrl *OrderController) Get(c *gin.Context) {
	orderId, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "invalid input", http.StatusUnprocessableEntity)
		return
	}

	var order domain.OrderTotal
	order, err = ctrl.service.Get(orderId)
	if err != nil {
		BadResponse(c, "no data found", http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "order retrieved", http.StatusOK, order)
}
