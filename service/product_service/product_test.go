package productservice_test

import (
	"errors"
	"fmt"
	"project/domain"
	"project/repository"
	productrepository "project/repository/product_repository"
	productservice "project/service/product_service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func base() (productservice.ProductService, *productrepository.ProductRepoMock) {
	log := *zap.NewNop()
	mockRepo := &productrepository.ProductRepoMock{}
	repo := repository.Repository{
		Product: mockRepo,
	}
	service := productservice.NewProductService(&repo, &log)

	return service, mockRepo
}

func TestShowAllProduct(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully show all products", func(t *testing.T) {
		page, limit := 1, 10

		mockProducts := []domain.Product{
			{ID: 1, Name: "Product A"},
			{ID: 2, Name: "Product B"},
		}
		mockRepo.On("ShowAllProduct", page, limit).
			Return(&mockProducts, 2, 1, nil).
			Once()

		products, count, totalPages, err := service.ShowAllProduct(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 2, count)
		assert.Equal(t, 1, totalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to show all products - Error from repository", func(t *testing.T) {
		page, limit := 1, 10

		mockRepo.On("ShowAllProduct", page, limit).
			Return(nil, 0, 0, fmt.Errorf("database error")).
			Once()

		products, count, totalPages, err := service.ShowAllProduct(page, limit)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, count)
		assert.Equal(t, 0, totalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to show all products - No products found", func(t *testing.T) {
		page, limit := 1, 10

		mockRepo.On("ShowAllProduct", page, limit).
			Return(nil, 0, 0, fmt.Errorf("no products found")).
			Once()

		products, count, totalPages, err := service.ShowAllProduct(page, limit)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, count)
		assert.Equal(t, 0, totalPages)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetProductByID(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully get product by ID", func(t *testing.T) {
		productID := 1

		mockProduct := &domain.Product{
			ID:   productID,
			Name: "Product A",
		}

		mockRepo.On("GetProductByID", productID).
			Return(mockProduct, nil).
			Once()

		product, err := service.GetProductByID(productID)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, productID, product.ID)
		assert.Equal(t, "Product A", product.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get product by ID - Product not found", func(t *testing.T) {
		productID := 2

		mockRepo.On("GetProductByID", productID).
			Return(nil, fmt.Errorf("product not found")).
			Once()

		product, err := service.GetProductByID(productID)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "product not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get product by ID - Database error", func(t *testing.T) {
		productID := 3

		mockRepo.On("GetProductByID", productID).
			Return(nil, fmt.Errorf("database error")).
			Once()

		product, err := service.GetProductByID(productID)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateProduct(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully create product", func(t *testing.T) {

		product := &domain.Product{
			Name:        "Product A",
			SKUProduct:  "SKU-2021",
			Price:       100,
			Description: "High-quality product",
		}

		mockRepo.On("CreateProduct", product).Return(nil).Once()

		err := service.CreateProduct(product)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Failed to create product - Repository error", func(t *testing.T) {

		product := &domain.Product{
			Name:        "Product B",
			SKUProduct:  "SKU-2022",
			Price:       150,
			Description: "Another high-quality product",
		}

		mockRepo.On("CreateProduct", product).Return(errors.New("failed to create product")).Once()

		err := service.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create product")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteProduct(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully delete product", func(t *testing.T) {
		productID := 1

		mockRepo.On("DeleteProduct", productID).Return(nil).Once()

		err := service.DeleteProduct(productID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to delete product - Repository error", func(t *testing.T) {
		productID := 2

		mockRepo.On("DeleteProduct", productID).Return(errors.New("failed to delete product")).Once()

		err := service.DeleteProduct(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete product")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to delete product - Product not found", func(t *testing.T) {
		productID := 3

		mockRepo.On("DeleteProduct", productID).Return(errors.New("product not found")).Once()

		err := service.DeleteProduct(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "product not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateProduct(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully update a product", func(t *testing.T) {
		productID := uint(1)
		product := &domain.Product{
			Name:        "Updated Product",
			SKUProduct:  "SKI-2022",
			Price:       150,
			Description: "Updated description",
		}

		mockRepo.On("UpdateProduct", productID, product).Return(nil)

		err := service.UpdateProduct(productID, product)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "UpdateProduct", productID, product)
	})

	t.Run("Failed to update product - No rows affected", func(t *testing.T) {
		productID := uint(2)
		product := &domain.Product{
			Name:        "Another Product",
			SKUProduct:  "SKI-2023",
			Price:       200,
			Description: "Another description",
		}

		mockRepo.On("UpdateProduct", productID, product).
			Return(fmt.Errorf("no record found with shipping_id %d", productID))

		err := service.UpdateProduct(productID, product)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("no record found with shipping_id %d", productID))
		mockRepo.AssertCalled(t, "UpdateProduct", productID, product)
	})

	t.Run("Failed to update product - Database error", func(t *testing.T) {
		productID := uint(3)
		product := &domain.Product{
			Name:        "Product with Error",
			SKUProduct:  "SKI-ERROR",
			Price:       300,
			Description: "Error description",
		}

		mockRepo.On("UpdateProduct", productID, product).
			Return(fmt.Errorf("database error"))

		err := service.UpdateProduct(productID, product)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertCalled(t, "UpdateProduct", productID, product)
	})
}
