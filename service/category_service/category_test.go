package categoryservice_test

import (
	"fmt"
	"project/domain"
	"project/repository"
	categoryrepositpry "project/repository/category_repositpry"
	categoryservice "project/service/category_service"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func base() (categoryservice.CategoryService, *categoryrepositpry.CategoryRepoMock) {
	log := *zap.NewNop()
	mockRepo := &categoryrepositpry.CategoryRepoMock{}
	repo := repository.Repository{
		Category: mockRepo,
	}
	service := categoryservice.NewCategoryService(&repo, &log)

	return service, mockRepo
}

func TestShowAllCategory(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully show all categories", func(t *testing.T) {
		page, limit := 1, 10
		mockCategories := []domain.Category{
			{ID: 1, Name: "Category A"},
			{ID: 2, Name: "Category B"},
		}

		mockRepo.On("ShowAllCategory", page, limit).
			Return(&mockCategories, 2, 1, nil).
			Once()

		categories, count, totalPages, err := service.ShowAllCategory(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		assert.Equal(t, 2, count)
		assert.Equal(t, 1, totalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to show all categories - Repository error", func(t *testing.T) {
		page, limit := 1, 10

		mockRepo.On("ShowAllCategory", page, limit).
			Return(nil, 0, 0, fmt.Errorf("database error")).
			Once()

		categories, count, totalPages, err := service.ShowAllCategory(page, limit)

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Equal(t, 0, count)
		assert.Equal(t, 0, totalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to show all categories - No categories found", func(t *testing.T) {
		page, limit := 1, 10

		mockRepo.On("ShowAllCategory", page, limit).
			Return(nil, 0, 0, fmt.Errorf("category not found")).
			Once()

		categories, count, totalPages, err := service.ShowAllCategory(page, limit)

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Equal(t, 0, count)
		assert.Equal(t, 0, totalPages)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully delete category", func(t *testing.T) {
		id := 1

		mockRepo.On("DeleteCategory", id).
			Return(nil).
			Once()

		err := service.DeleteCategory(id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to delete category - Category not found", func(t *testing.T) {
		id := 999 // Assume this ID does not exist

		mockRepo.On("DeleteCategory", id).
			Return(fmt.Errorf("category not found")).
			Once()

		err := service.DeleteCategory(id)

		assert.Error(t, err)
		assert.EqualError(t, err, "category not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to delete category - Repository error", func(t *testing.T) {
		id := 1

		mockRepo.On("DeleteCategory", id).
			Return(fmt.Errorf("database error")).
			Once()

		err := service.DeleteCategory(id)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCategoryByID(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully get category by ID", func(t *testing.T) {
		id := 1
		mockCategory := &domain.Category{
			ID:   1,
			Name: "Electronics",
		}

		mockRepo.On("GetCategoryByID", id).
			Return(mockCategory, nil).
			Once()

		category, err := service.GetCategoryByID(id)

		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, mockCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get category by ID - Category not found", func(t *testing.T) {
		id := 999 // Assume this ID does not exist

		mockRepo.On("GetCategoryByID", id).
			Return(nil, fmt.Errorf("category not found")).
			Once()

		category, err := service.GetCategoryByID(id)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.EqualError(t, err, "category not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to get category by ID - Repository error", func(t *testing.T) {
		id := 1

		mockRepo.On("GetCategoryByID", id).
			Return(nil, fmt.Errorf("database error")).
			Once()

		category, err := service.GetCategoryByID(id)

		assert.Error(t, err)
		assert.Nil(t, category)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateCategory(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully update category", func(t *testing.T) {
		id := 1
		updateData := &domain.Category{
			Name: "Updated Category Name",
		}

		mockRepo.On("UpdateCategory", id, updateData).
			Return(nil).
			Once()

		err := service.UpdateCategory(id, updateData)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to update category - Category not found", func(t *testing.T) {
		id := 999 // Assume this ID does not exist
		updateData := &domain.Category{
			Name: "Nonexistent Category",
		}

		mockRepo.On("UpdateCategory", id, updateData).
			Return(fmt.Errorf("category not found")).
			Once()

		err := service.UpdateCategory(id, updateData)

		assert.Error(t, err)
		assert.EqualError(t, err, "category not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to update category - Repository error", func(t *testing.T) {
		id := 1
		updateData := &domain.Category{
			Name: "Category with DB Error",
		}

		mockRepo.On("UpdateCategory", id, updateData).
			Return(fmt.Errorf("database error")).
			Once()

		err := service.UpdateCategory(id, updateData)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateCategory(t *testing.T) {
	service, mockRepo := base()

	t.Run("Successfully create category", func(t *testing.T) {
		newCategory := &domain.Category{
			Name: "New Category",
		}

		mockRepo.On("CreateCategory", newCategory).
			Return(nil).
			Once()

		err := service.CreateCategory(newCategory)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed to create category - Repository error", func(t *testing.T) {
		newCategory := &domain.Category{
			Name: "Invalid Category",
		}

		mockRepo.On("CreateCategory", newCategory).
			Return(fmt.Errorf("database error")).
			Once()

		err := service.CreateCategory(newCategory)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}
