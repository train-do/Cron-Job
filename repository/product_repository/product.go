package productrepository

import (
	"fmt"
	"log"
	"math"
	"project/domain"
	"project/helper"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepo interface {
	ShowAllProduct(page, limit int) (*[]domain.Product, int, int, error)
	GetProductByID(id int) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	DeleteProduct(id int) error
	UpdateProduct(productID uint, product *domain.Product) error
}

type productRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewProductRepo(db *gorm.DB, log *zap.Logger) ProductRepo {
	return &productRepo{db, log}
}

func (pr *productRepo) ShowAllProduct(page, limit int) (*[]domain.Product, int, int, error) {
	pr.log.Info("Fetching all products", zap.Int("page", page), zap.Int("limit", limit))

	productList := []domain.Product{}
	var count int64

	if err := pr.db.Model(&domain.Product{}).Count(&count).Error; err != nil {
		pr.log.Error("Error counting products", zap.Error(err))
		return nil, 0, 0, err
	}

	result := pr.db.Scopes(helper.Paginate(uint(page), uint(limit))).
		Preload("ProductVariant").
		Preload("Image").
		Find(&productList)

	if result.Error != nil {
		pr.log.Error("Error fetching products", zap.Error(result.Error))
		return nil, 0, 0, result.Error
	}

	if result.RowsAffected == 0 {
		pr.log.Warn("No products found")
		return nil, 0, 0, fmt.Errorf("no products found")
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	pr.log.Info("Successfully fetched products", zap.Int("totalCount", int(count)), zap.Int("totalPages", totalPages))
	return &productList, int(count), totalPages, nil
}

func (pr *productRepo) GetProductByID(id int) (*domain.Product, error) {
	pr.log.Info("Fetching product by ID", zap.Int("id", id))

	product := domain.Product{}

	result := pr.db.Model(&product).Where("id = ?", id).
		Preload("ProductVariant").
		Preload("Image").First(&product)

	if result.Error != nil {
		pr.log.Error("Error fetching product", zap.Int("id", id), zap.Error(result.Error))
		return nil, fmt.Errorf("product not found")
	}

	pr.log.Info("Successfully fetched product", zap.Int("id", id))
	return &product, nil
}

func (pr *productRepo) CreateProduct(product *domain.Product) error {
	pr.log.Info("Starting product creation", zap.String("productName", product.Name))

	var wg sync.WaitGroup
	var err error

	err = pr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			pr.log.Error("Failed to create product", zap.String("productName", product.Name), zap.Error(err))
			return fmt.Errorf("failed to create product: %w", err)
		}

		if err := tx.Create(&domain.ProductVariant{
			ProductID: product.ID,
		}).Error; err != nil {
			pr.log.Error("Failed to create product variant", zap.Error(err))
			return fmt.Errorf("failed to create product variant: %w", err)
		}

		if err := tx.Create(&domain.Image{
			ProductID: product.ID,
		}).Error; err != nil {
			pr.log.Error("Failed to create image", zap.Error(err))
			log.Printf("failed to create image: %v", err)
			return fmt.Errorf("failed to create image: %w", err)
		}

		wg.Wait()

		if err != nil {
			pr.log.Error("Error occurred during transaction", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		pr.log.Error("Transaction failed", zap.String("productName", product.Name), zap.Error(err))
		log.Printf("Transaction failed: %v", err)
	} else {
		pr.log.Info("Product creation successful", zap.String("productName", product.Name))
	}

	return err
}

func (pr *productRepo) UpdateProduct(productID uint, product *domain.Product) error {
	pr.log.Info("Updating product", zap.Uint("productID", productID), zap.String("productName", product.Name))

	result := pr.db.Model(&product).
		Where("id = ?", productID).Updates(product)

	if result.Error != nil {
		pr.log.Error("Failed to update product", zap.Uint("productID", productID), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		pr.log.Warn("No record found to update", zap.Uint("productID", productID))
		return fmt.Errorf("no record found with shipping_id %d", productID)
	}

	pr.log.Info("Successfully updated product", zap.Uint("productID", productID))
	return nil
}

func (pr *productRepo) DeleteProduct(id int) error {

	pr.log.Info("Deleting product", zap.Int("productID", id))

	result := pr.db.Delete(&domain.Product{}, id)
	if result.Error != nil {
		pr.log.Error("Failed to delete product", zap.Int("productID", id), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		pr.log.Warn("Product not found for deletion", zap.Int("productID", id))
		return fmt.Errorf("product not found")
	}

	pr.log.Info("Successfully deleted product", zap.Int("productID", id))
	return nil
}
