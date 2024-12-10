package categoryrepositpry

import (
	"project/domain"

	"github.com/stretchr/testify/mock"
)

type CategoryRepoMock struct {
	mock.Mock
}

func (cr *CategoryRepoMock) ShowAllCategory(page, limit int) (*[]domain.Category, int, int, error) {
	args := cr.Called(page, limit)
	if categories, ok := args.Get(0).(*[]domain.Category); ok {
		return categories, args.Int(1), args.Int(2), args.Error(3)
	}
	return nil, 0, 0, args.Error(3)
}

func (cr *CategoryRepoMock) DeleteCategory(id int) error {
	args := cr.Called(id)
	return args.Error(0)
}

func (cr *CategoryRepoMock) GetCategoryByID(id int) (*domain.Category, error) {
	args := cr.Called(id)
	if category, ok := args.Get(0).(*domain.Category); ok {
		return category, args.Error(1)
	}
	return nil, args.Error(1)
}

func (cr *CategoryRepoMock) CreateCategory(category *domain.Category) error {
	args := cr.Called(category)
	if err, ok := args.Error(0).(error); ok {
		return err
	}
	return nil
}

func (cr *CategoryRepoMock) UpdateCategory(id int, category *domain.Category) error {
	args := cr.Called(id, category)
	if err, ok := args.Error(0).(error); ok {
		return err
	}
	return nil
}
