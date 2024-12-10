package productrepository

import (
	"project/domain"

	"github.com/stretchr/testify/mock"
)

type ProductRepoMock struct {
	mock.Mock
}

func (pr *ProductRepoMock) ShowAllProduct(page, limit int) (*[]domain.Product, int, int, error) {
	args := pr.Called(page, limit)
	if products := args.Get(0); products != nil {
		return products.(*[]domain.Product), args.Int(1), args.Int(2), args.Error(3)
	}
	return nil, 0, 0, args.Error(3)
}

func (pr *ProductRepoMock) GetProductByID(id int) (*domain.Product, error) {
	args := pr.Called(id)
	if product, ok := args.Get(0).(*domain.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

func (pr *ProductRepoMock) CreateProduct(product *domain.Product) error {
	args := pr.Called(product)
	return args.Error(0)
}

func (pr *ProductRepoMock) DeleteProduct(id int) error {
	args := pr.Called(id)
	return args.Error(0)
}

func (pr *ProductRepoMock) UpdateProduct(productID uint, product *domain.Product) error {
	args := pr.Called(productID, product)
	return args.Error(0)
}
