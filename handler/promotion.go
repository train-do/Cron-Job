package handler

import (
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerPromotion struct {
	service service.ServicePromotion
	logger  *zap.Logger
}

func NewControllerPromotion(service service.ServicePromotion, logger *zap.Logger) *ControllerPromotion {
	return &ControllerPromotion{service: service, logger: logger}
}

// @Summary Get All Promotion
// @Description Endpoint Fetch All Promotion
// @Tags promotions
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response{data=[]domain.Promotion} "Get All Success"
// @Failure 500 {object} handler.Response "server error"
// @Router  /promotion [get]
func (ctrl *ControllerPromotion) GetAll(c *gin.Context) {
	banners, err := ctrl.service.GetAll()
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "Get Promotions success", http.StatusOK, banners)
}

// @Summary Get a promotion by ID
// @Description Get details of a specific promotion using the provided ID
// @Tags promotions
// @Accept  json
// @Produce  json
// @Param id path int true "Promotion ID"
// @Success 200 {object} handler.Response{data=domain.Promotion}  "Promotion details"
// @Failure 400 {object} handler.Response  "Bad Request"
// @Failure 500 {object} handler.Response  "Internal server error"
// @Router /promotion/{id} [get]
func (ctrl *ControllerPromotion) GetById(c *gin.Context) {
	id, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	banner, err := ctrl.service.GetById(id)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "Get Promotion success", http.StatusOK, banner)
}

// @Summary Create a new promotion
// @Description Create a new promotion by sending the promotion data in the request body
// @Tags promotions
// @Accept  json
// @Produce  json
// @Param promotion body domain.Promotion true "Promotion data"
// @Success 200 {object} handler.Response{data=domain.Promotion}  "Promotion details"
// @Failure 400 {object} handler.Response  "Bad Request"
// @Failure 500 {object} handler.Response  "Internal server error"
// @Router /promotion [post]
func (ctrl *ControllerPromotion) Create(c *gin.Context) {
	var data domain.Promotion
	if err := c.ShouldBind(&data); err != nil {
		BadResponse(c, "Bad Request (Body)", http.StatusBadRequest)
		return
	}
	err := ctrl.service.Create(&data)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "Create Promotion success", http.StatusCreated, data)
}

// @Summary Delete a promotion by ID
// @Description Delete a promotion by its ID from the database
// @Tags promotions
// @Accept  json
// @Produce  json
// @Param id path int true "Promotion ID"
// @Success 200 {object} handler.Response{data=domain.Promotion}  "Promotion details"
// @Failure 400 {object} handler.Response  "Bad Request"
// @Failure 500 {object} handler.Response  "Internal server error"
// @Router /promotion/{id} [delete]
func (ctrl *ControllerPromotion) Delete(c *gin.Context) {
	id, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	var data domain.Promotion
	data.ID = id
	err = ctrl.service.Delete(&data)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "Delete Promotion success", http.StatusOK, data)
}
