package productservice

import (
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type ProductService interface {
	ShowAllProduct(page, limit int) (*[]domain.Product, int, int, error)
	GetProductByID(id int) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	DeleteProduct(id int) error
	UpdateProduct(productID uint, product *domain.Product) error
}

type productService struct {
	repo *repository.Repository
	log  *zap.Logger
}

func NewProductService(repo *repository.Repository, log *zap.Logger) ProductService {
	return &productService{repo, log}
}

func (ps *productService) ShowAllProduct(page, limit int) (*[]domain.Product, int, int, error) {
	ps.log.Info("Fetching all products", zap.Int("page", page), zap.Int("limit", limit))

	products, count, totalPages, err := ps.repo.Product.ShowAllProduct(page, limit)
	if err != nil {
		ps.log.Error("Error fetching products", zap.Error(err))
		return nil, 0, 0, err
	}

	ps.log.Info("Successfully fetched products", zap.Int("count", count), zap.Int("totalPages", totalPages))
	return products, count, totalPages, nil
}

func (ps *productService) GetProductByID(id int) (*domain.Product, error) {
	ps.log.Info("Fetching product by ID", zap.Int("id", id))

	product, err := ps.repo.Product.GetProductByID(id)
	if err != nil {
		ps.log.Error("Error fetching product", zap.Error(err))
		return nil, err
	}

	ps.log.Info("Successfully fetched product", zap.Int("id", id))
	return product, nil
}

func (ps *productService) CreateProduct(product *domain.Product) error {
	ps.log.Info("Creating product", zap.String("name", product.Name))

	err := ps.repo.Product.CreateProduct(product)
	if err != nil {
		ps.log.Error("Error creating product", zap.Error(err))
		return err
	}

	ps.log.Info("Successfully created product", zap.String("name", product.Name))
	return nil
}

func (ps *productService) DeleteProduct(id int) error {
	ps.log.Info("Deleting product", zap.Int("id", id))

	err := ps.repo.Product.DeleteProduct(id)
	if err != nil {
		ps.log.Error("Error deleting product", zap.Error(err))
		return err
	}

	ps.log.Info("Successfully deleted product", zap.Int("id", id))
	return nil
}

func (ps *productService) UpdateProduct(productID uint, product *domain.Product) error {
	ps.log.Info("Updating product", zap.Uint("productID", productID), zap.String("name", product.Name))

	if err := ps.repo.Product.UpdateProduct(productID, product); err != nil {
		ps.log.Error("Error updating product", zap.Error(err))
		return err
	}

	ps.log.Info("Successfully updated product", zap.Uint("productID", productID))
	return nil
}
