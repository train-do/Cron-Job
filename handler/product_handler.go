package handler

import (
	"encoding/json"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductHandler interface {
	ShowAllProduct(c *gin.Context)
	GetProductByID(c *gin.Context)
	CreateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
}

type productHandler struct {
	service *service.Service
	log     *zap.Logger
}

func NewProductHandler(service *service.Service, log *zap.Logger) ProductHandler {
	return &productHandler{service, log}
}

// @Summary Get all products with pagination
// @Description Fetches a paginated list of all products
// @Tags Product
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} handler.Response{data=[]domain.Product} "Successfully retrieved products"
// @Failure 400 {object} handler.Response "Invalid query parameters"
// @Failure 404 {object} handler.Response "Products not found"
// @Failure 500 {object} handler.Response "Internal server error"
// @Router /products [get]
func (ph *productHandler) ShowAllProduct(c *gin.Context) {
	ph.log.Info("Fetching all products", zap.String("queryPage", c.Query("page")), zap.String("queryLimit", c.Query("limit")))

	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
		ph.log.Warn("Invalid page number, defaulting to 1")
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 10 {
		limit = 10
		ph.log.Warn("Limit too low, defaulting to 10")
	}

	products, count, totalPages, err := ph.service.Product.ShowAllProduct(page, limit)
	if err != nil {
		ph.log.Error("Failed to retrieve products", zap.Int("page", page), zap.Int("limit", limit), zap.Error(err))
		BadResponse(c, "Product Not Found", http.StatusNotFound)
		return
	}

	ph.log.Info("Successfully retrieved products", zap.Int("page", page), zap.Int("limit", limit), zap.Int("count", count), zap.Int("totalPages", totalPages))
	GoodResponseWithPage(c, "Successfully Retrieved Products", http.StatusOK, count, totalPages, page, limit, products)
}

// Get Product By ID
// @Summary Get Product By ID
// @Description Get Product By ID
// @Tags Products
// @Param id path int true "Product ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "Successfully Retrieved Product"
// @Failure 404 {object} handler.Response "Product Not Found"
// @Router  /products/{id} [get]
func (ph *productHandler) GetProductByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ph.log.Info("Fetching product by ID", zap.Int("productID", id))

	product, err := ph.service.Product.GetProductByID(id)
	if err != nil {
		ph.log.Error("Product not found", zap.Int("productID", id), zap.Error(err))
		BadResponse(c, "Product Not Found", http.StatusNotFound)
		return
	}

	ph.log.Info("Successfully retrieved product", zap.Int("productID", id))
	GoodResponseWithData(c, "Successfully Retrieved Product", http.StatusOK, product)
}

// Create Product
// @Summary Create Product
// @Description Create a new product with variants and images
// @Tags Products
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "Product Name"
// @Param sku_product formData string true "Product SKU"
// @Param price formData int true "Product Price"
// @Param description formData string true "Product Description"
// @Param images formData file true "Product Images" multiple
// @Param variants formData string true "Product Variants in JSON format"
// @Success 201 {object} handler.Response{data=domain.Product} "Product created successfully"
// @Failure 400 {object} handler.Response "Invalid form data"
// @Failure 500 {object} handler.Response "Failed to create product"
// @Router /products [post]
func (ph *productHandler) CreateProduct(c *gin.Context) {
	ph.log.Info("Starting product creation")

	form, err := c.MultipartForm()
	if err != nil {
		ph.log.Error("Error reading form data", zap.Error(err))
		BadResponse(c, "Invalid form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := form.File["images"]
	for _, file := range files {
		ph.log.Info("Processing uploaded file", zap.String("fileName", file.Filename), zap.Int64("fileSize", file.Size))
	}

	var wg sync.WaitGroup
	responses, err := helper.Upload(&wg, files)
	if err != nil {
		ph.log.Error("Failed to upload images", zap.Error(err))
		BadResponse(c, "Failed to upload images: "+err.Error(), http.StatusInternalServerError)
		return
	}

	name := c.PostForm("name")
	skuProduct := c.PostForm("sku_product")
	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		ph.log.Error("Invalid price value", zap.String("price", c.PostForm("price")), zap.Error(err))
		BadResponse(c, "Invalid price value", http.StatusBadRequest)
		return
	}

	ph.log.Info("Parsed product data", zap.String("name", name), zap.String("skuProduct", skuProduct), zap.Int("price", price))

	var images []*domain.Image
	for _, response := range responses {
		images = append(images, &domain.Image{
			URLPath: response.Data.Url,
		})
	}

	variantData := c.PostForm("variants")
	ph.log.Info("Parsing variant data", zap.String("variantData", variantData))

	var productVariants []*domain.ProductVariant
	err = json.Unmarshal([]byte(variantData), &productVariants)
	if err != nil {
		ph.log.Error("Invalid variant data", zap.String("variantData", variantData), zap.Error(err))
		BadResponse(c, "Invalid variant data: "+err.Error(), http.StatusBadRequest)
		return
	}

	product := domain.Product{
		Name:           name,
		SKUProduct:     skuProduct,
		Price:          float64(price),
		Description:    c.PostForm("description"),
		Image:          images,
		ProductVariant: productVariants,
	}

	ph.log.Info("Creating product", zap.String("productName", product.Name), zap.Float64("price", product.Price))

	if err := ph.service.Product.CreateProduct(&product); err != nil {
		ph.log.Error("Failed to create product", zap.String("productName", product.Name), zap.Error(err))
		BadResponse(c, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ph.log.Info("Product created successfully", zap.String("productName", product.Name))
	GoodResponseWithData(c, "Product created successfully", http.StatusCreated, product)
}

// Delete Product
// @Summary Delete Product
// @Description Delete Product
// @Tags Products
// @Accept  json
// @Param id path int true "Product ID"
// @Produce  json
// @Success 200 {object} handler.Response "Product Deleted successfully"
// @Failure 404 {object} handler.Response "Failed to Delete product"
// @Router  /products/{id} [delete]
func (ph *productHandler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ph.log.Info("Attempting to delete product", zap.Int("productID", id))

	if err := ph.service.Product.DeleteProduct(id); err != nil {
		ph.log.Error("Failed to delete product", zap.Int("productID", id), zap.Error(err))
		BadResponse(c, "Failed to Delete product: "+err.Error(), http.StatusNotFound)
		return
	}

	ph.log.Info("Product deleted successfully", zap.Int("productID", id))
	GoodResponseWithData(c, "Product Deleted successfully", http.StatusOK, id)
}

// Update Product
// @Summary Update Product
// @Description Update the details of a product
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Param name body string true "Product Name"
// @Param sku_product body string true "Product SKU"
// @Param price body int true "Product Price"
// @Param description body string true "Product Description"
// @Success 200 {object} handler.Response "Product Updated successfully"
// @Failure 400 {object} handler.Response "Failed to Update product"
// @Failure 500 {object} handler.Response "Invalid Payload Request"
// @Router /products/{id} [put]
func (ph *productHandler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ph.log.Info("Attempting to update product", zap.Int("productID", id))

	product := domain.Product{}
	if err := c.ShouldBindJSON(&product); err != nil {
		ph.log.Error("Failed to bind product data", zap.Int("productID", id), zap.Error(err))
		BadResponse(c, "Failed to Update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ph.log.Info("Product data bound successfully", zap.Int("productID", id), zap.String("productName", product.Name))

	if err := ph.service.Product.UpdateProduct(uint(id), &product); err != nil {
		ph.log.Error("Failed to update product", zap.Int("productID", id), zap.Error(err))
		BadResponse(c, "Failed to Update product: "+err.Error(), http.StatusBadRequest)
		return
	}

	ph.log.Info("Product updated successfully", zap.Int("productID", id), zap.String("productName", product.Name))
	GoodResponseWithData(c, "Product Updated successfully", http.StatusOK, product)
}
