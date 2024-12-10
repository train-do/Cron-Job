package handler

import (
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerStock struct {
	service service.ServiceStock
	log     *zap.Logger
}

func NewServiceStock(service service.ServiceStock, log *zap.Logger) *ControllerStock {
	return &ControllerStock{service, log}
}

// @Summary Get Stock Details
// @Description Get details of the stock by product variant ID.
// @Tags Stock
// @Accept json
// @Produce json
// @Param productVariantId path int true "Product Variant ID"
// @Success 200 {object} handler.Response{data=domain.ResponseStock} "Stock details retrieved successfully"
// @Failure 400 {object} handler.Response "Invalid parameters or bad request"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /stock/{productVariantId} [get]
func (ctrl *ControllerStock) GetDetails(c *gin.Context) {
	id, err := helper.Uint(c.Param("productVariantId"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	data, err := ctrl.service.GetDetails(int(id))
	if err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "Get Detail Stock success", http.StatusOK, data)
}

type FormStock struct {
	NewStock int
}

// @Summary Edit Stock Details
// @Description Edit the stock quantity by product variant ID.
// @Tags Stock
// @Accept json
// @Produce json
// @Param productVariantId path int true "Product Variant ID"
// @Param body body FormStock true "New Stock Quantity"
// @Success 200 {object} handler.Response{data=domain.Stock} "Stock updated successfully"
// @Failure 400 {object} handler.Response "Invalid parameters or bad request"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /stock/{productVariantId} [put]
func (ctrl *ControllerStock) Edit(c *gin.Context) {
	id, err := helper.Uint(c.Param("productVariantId"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	newStock := FormStock{}
	if err := c.ShouldBindJSON(&newStock); err != nil {
		BadResponse(c, "Bad Request (Body)", http.StatusBadRequest)
		return
	}
	// fmt.Println(c.PostForm("newStock"), "<<<<<<<")
	// newStock, err := helper.Uint(c.PostForm("newStock"))
	// if err != nil {
	// 	BadResponse(c, "Bad Request (Body)", http.StatusBadRequest)
	// 	return
	// }
	data, err := ctrl.service.Edit(int(id), newStock.NewStock)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "Edit Stock success", http.StatusOK, data)
}

// @Summary Delete Stock
// @Description Delete stock by product variant ID.
// @Tags Stock
// @Accept json
// @Produce json
// @Param id path int true "Product Variant ID"
// @Success 200 {object} handler.Response{data=domain.Stock} "Stock deleted successfully"
// @Failure 400 {object} handler.Response "Invalid parameters or bad request"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /stock/{id} [delete]
func (ctrl *ControllerStock) Delete(c *gin.Context) {
	id, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	var data domain.Stock
	data.ID = id
	if err := ctrl.service.Delete(&data); err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "Delete History Stock success", http.StatusOK, data)
}
