package handler

import (
	"log"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControllerBanner struct {
	service service.ServiceBanner
	logger  *zap.Logger
}

func NewControllerBanner(service service.ServiceBanner, logger *zap.Logger) *ControllerBanner {
	return &ControllerBanner{service: service, logger: logger}
}

// @Summary Get All Banner
// @Description Endpoint Fetch All Banner
// @Tags Banner
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response{data=[]domain.Banner} "Get All Success"
// @Failure 500 {object} handler.Response "server error"
// @Router  /banner [get]
func (ctrl *ControllerBanner) GetAll(c *gin.Context) {
	banners, err := ctrl.service.GetAll()
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "Get Banners success", http.StatusOK, banners)
}

// @Summary Get Banner by ID
// @Description Get a banner details by its ID.
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} handler.Response{data=domain.Banner} "Success"
// @Failure 400 {object} handler.Response "Bad Request"
// @Failure 404 {object} handler.Response "Banner Not Found"
// @Router /banner/{id} [get]
func (ctrl *ControllerBanner) GetById(c *gin.Context) {
	id, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	banner, err := ctrl.service.GetById(id)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "Get Banner success", http.StatusOK, banner)
}

// @Summary Create a new Banner
// @Description Create a new banner with a title, path, start date, end date, and image upload.
// @Tags Banner
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Banner Title"
// @Param pathPage formData string true "Path Page"
// @Param startDate formData string true "Start Date (yyyy-mm-dd)"
// @Param endDate formData string true "End Date (yyyy-mm-dd)"
// @Param images formData file false "Banner Image"
// @Success 201 {string} handler.Response{data=domain.Banner} "Banner successfully created"
// @Failure 400 {object} handler.Response "Invalid form data"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /banner [post]
func (ctrl *ControllerBanner) Create(c *gin.Context) {
	doUpload := true
	var respThirdParty []domain.CdnResponse
	_, err := c.FormFile("images")
	if err != nil {
		doUpload = false
		respThirdParty = append(respThirdParty, domain.CdnResponse{Data: struct {
			FileId      string "json:\"fileId\""
			Name        string "json:\"name\""
			Size        int    "json:\"size\""
			VersionInfo struct {
				Id   string "json:\"id\""
				Name string "json:\"name\""
			} "json:\"versionInfo\""
			FilePath     string      "json:\"filePath\""
			Url          string      "json:\"url\""
			FileType     string      "json:\"fileType\""
			Height       int         "json:\"height\""
			Width        int         "json:\"width\""
			ThumbnailUrl string      "json:\"thumbnailUrl\""
			AITags       interface{} "json:\"AITags\""
		}{Url: ""}})
	}
	if doUpload {
		form, err := c.MultipartForm()
		if err != nil {
			log.Println("Error reading form data:", err)
			BadResponse(c, "Invalid form data: "+err.Error(), http.StatusBadRequest)
			return
		}
		files := form.File["images"]
		for _, file := range files {
			log.Println("File size:", file.Size)
		}

		var wg sync.WaitGroup
		respThirdParty, err = helper.Upload(&wg, files)
		if err != nil {
			BadResponse(c, "Failed to upload images: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	banner := domain.Banner{
		Title:     c.PostForm("title"),
		PathPage:  c.PostForm("pathPage"),
		StartDate: c.PostForm("startDate"),
		EndDate:   c.PostForm("endDate"),
		IsPublish: false,
		ImageUrl:  respThirdParty[0].Data.Url,
	}
	err = ctrl.service.Create(&banner)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "This Banner was successfully added", http.StatusCreated, banner)
}

// @Summary Edit a new Banner
// @Description Edit a new banner with a title, path, start date, end date, and image upload.
// @Tags Banner
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Banner Title"
// @Param pathPage formData string true "Path Page"
// @Param startDate formData string true "Start Date (yyyy-mm-dd)"
// @Param endDate formData string true "End Date (yyyy-mm-dd)"
// @Param images formData file false "Banner Image"
// @Success 201 {string} handler.Response{data=domain.Banner} "Banner successfully created"
// @Failure 400 {object} handler.Response "Invalid form data"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /banner/{id} [put]
func (ctrl *ControllerBanner) Edit(c *gin.Context) {
	id, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	doUpload := true
	var respThirdParty []domain.CdnResponse
	_, err = c.FormFile("images")
	if err != nil {
		doUpload = false
		respThirdParty = append(respThirdParty, domain.CdnResponse{Data: struct {
			FileId      string "json:\"fileId\""
			Name        string "json:\"name\""
			Size        int    "json:\"size\""
			VersionInfo struct {
				Id   string "json:\"id\""
				Name string "json:\"name\""
			} "json:\"versionInfo\""
			FilePath     string      "json:\"filePath\""
			Url          string      "json:\"url\""
			FileType     string      "json:\"fileType\""
			Height       int         "json:\"height\""
			Width        int         "json:\"width\""
			ThumbnailUrl string      "json:\"thumbnailUrl\""
			AITags       interface{} "json:\"AITags\""
		}{Url: ""}})
	}
	if doUpload {
		form, err := c.MultipartForm()
		if err != nil {
			log.Println("Error reading form data:", err)
			BadResponse(c, "Invalid form data: "+err.Error(), http.StatusBadRequest)
			return
		}
		files := form.File["images"]
		for _, file := range files {
			log.Println("File size:", file.Size)
		}

		var wg sync.WaitGroup
		respThirdParty, err = helper.Upload(&wg, files)
		if err != nil {
			BadResponse(c, "Failed to upload images: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	banner := domain.Banner{
		ID:        id,
		Title:     c.PostForm("title"),
		PathPage:  c.PostForm("pathPage"),
		StartDate: c.PostForm("startDate"),
		EndDate:   c.PostForm("endDate"),
		IsPublish: strings.ToLower(c.PostForm("isPublish")) == "true",
		ImageUrl:  respThirdParty[0].Data.Url,
	}
	err = ctrl.service.Edit(&banner)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "This Banner was successfully updated", http.StatusCreated, banner)
}

// @Summary Delete a Banner
// @Description Delete a banner by its unique ID.
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} handler.Response{data=domain.Banner} "Banner successfully deleted"
// @Failure 400 {object} handler.Response "Invalid parameters or bad request"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /banner/{id} [delete]
func (ctrl *ControllerBanner) Delete(c *gin.Context) {
	id, err := helper.Uint(c.Param("id"))
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	var banner domain.Banner
	banner.ID = id
	err = ctrl.service.Delete(&banner)
	if err != nil {
		BadResponse(c, "Bad Request (Params)", http.StatusBadRequest)
		return
	}
	GoodResponseWithData(c, "This Banner was successfully deleted", http.StatusCreated, banner)
}
