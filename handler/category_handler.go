package handler

import (
	"mime/multipart"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler interface {
	ShowAllCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	CreateCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
}

type categoryHandler struct {
	log     *zap.Logger
	service *service.Service
}

func NewCategoryHandler(log *zap.Logger, service *service.Service) CategoryHandler {
	return &categoryHandler{log, service}
}

// @Summary Show all categories
// @Description Retrieves all categories with pagination support
// @Tags Category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} handler.Response{data=[]domain.Category} "Successfully retrieved categories"
// @Failure 404 {object} handler.Response "Failed to retrieve categories"
// @Router /category [get]
func (ch *categoryHandler) ShowAllCategory(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
		ch.log.Warn("Invalid page number, defaulting to 1")
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 10 {
		limit = 10
		ch.log.Warn("Limit too low, defaulting to 10")
	}

	categories, count, totalPages, err := ch.service.Category.ShowAllCategory(page, limit)
	if err != nil {
		BadResponse(c, "Failed to retrived categories: "+err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithPage(c, "successfully retrived categories", http.StatusOK, count, totalPages, page, limit, categories)
}

// @Summary Delete a category
// @Description Deletes a category by its ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} handler.Response{data=domain.Category} "Successfully deleted category"
// @Failure 404 {object} handler.Response "Failed to delete category"
// @Router /category/{id} [delete]
func (ch *categoryHandler) DeleteCategory(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	if err := ch.service.Category.DeleteCategory(id); err != nil {
		BadResponse(c, "Failed to deleted categoriy: "+err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "successfully deleted category", http.StatusOK, id)

}

// @Summary Get a category by ID
// @Description Retrieves a category by its ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} handler.Response{data=domain.Category} "Successfully retrieved category"
// @Failure 404 {object} handler.Response "Failed to retrieve category"
// @Router /category/{id} [get]
func (ch *categoryHandler) GetCategoryByID(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	category, err := ch.service.Category.GetCategoryByID(id)
	if err != nil {
		BadResponse(c, "Failed to Retrieved categoriy: "+err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "successfully Retrieved category", http.StatusOK, category)
}

// @Summary Create a new category
// @Description Creates a new category with an image and name
// @Tags Category
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Category Name"
// @Param images formData file true "Category Image"
// @Success 201 {object} handler.Response{data=domain.Category} "Category created successfully"
// @Failure 400 {object} handler.Response "Bad request, invalid data"
// @Failure 500 {object} handler.Response "Internal server error"
// @Router /category [post]
func (ch *categoryHandler) CreateCategory(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		BadResponse(c, "Invalid form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		BadResponse(c, "No image provided", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup

	responses, err := helper.Upload(&wg, []*multipart.FileHeader{files[0]})
	if err != nil || len(responses) == 0 {
		BadResponse(c, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	name := c.PostForm("name")
	if name == "" {
		BadResponse(c, "Name is required", http.StatusBadRequest)
		return
	}

	imageURL := responses[0].Data.Url

	// Buat entitas kategori baru
	category := domain.Category{
		Name: name,
		Icon: imageURL,
	}

	// Simpan kategori menggunakan service
	err = ch.service.Category.CreateCategory(&category)
	if err != nil {
		BadResponse(c, "Failed to create category: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Berikan respon sukses
	GoodResponseWithData(c, "Category created successfully", http.StatusCreated, category)
}

// @Summary Update an existing category
// @Description Updates a category
// @Tags Category
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Category ID"
// @Param name formData string false "Category Name"
// @Param images formData file false "Category Image"
// @Success 200 {object} handler.Response{data=domain.Category} "Category updated successfully"
// @Failure 400 {object} handler.Response "Bad request, invalid data"
// @Failure 404 {object} handler.Response "Category not found"
// @Failure 500 {object} handler.Response "Internal server error"
// @Router /category/{id} [put]
func (ch *categoryHandler) UpdateCategory(c *gin.Context) {

	category := domain.Category{}
	id, _ := strconv.Atoi(c.Param("id"))

	form, _ := c.MultipartForm()

	var imageURL string
	files := form.File["images"]
	if len(files) > 0 {
		var wg sync.WaitGroup
		responses, err := helper.Upload(&wg, []*multipart.FileHeader{files[0]})
		if err != nil || len(responses) == 0 {
			BadResponse(c, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		imageURL = responses[0].Data.Url
	}

	name := c.PostForm("name")
	if name == "" {
		name = category.Name
	}

	category = domain.Category{
		ID:   uint(id),
		Name: name,
		Icon: imageURL,
	}

	err := ch.service.Category.UpdateCategory(id, &category)
	if err != nil {
		BadResponse(c, "Category not found: "+err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "Category updated successfully", http.StatusOK, category)
}
