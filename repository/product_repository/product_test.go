package productrepository_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"project/domain"
	"project/helper"
	productrepository "project/repository/product_repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestShowAllProduct(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	productRepo := productrepository.NewProductRepo(db, &log)

	t.Run("Successfully show all products", func(t *testing.T) {
		page, limit := 1, 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" WHERE "products"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		// Mock data query with pagination
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL LIMIT $1`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "Product A").
				AddRow(2, "Product B"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."product_id" IN ($1,$2) AND "images"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).
				AddRow(1, 1).
				AddRow(2, 2))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product_variants" WHERE "product_variants"."product_id" IN ($1,$2) `)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).
				AddRow(1, 1).
				AddRow(2, 2))

		products, totalCount, totalPages, err := productRepo.ShowAllProduct(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, *products, 2)

		assert.Equal(t, "Product A", (*products)[0].Name)
		assert.Equal(t, "Product B", (*products)[1].Name)

		assert.Equal(t, 5, totalCount)
		assert.Equal(t, 3, totalPages) // 5 products / 2 per page = 3 pages
	})

	t.Run("Failed to show all products due to database error", func(t *testing.T) {
		page, limit := 1, 2

		// Mock count query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" WHERE "products"."deleted_at" IS NULL`)).
			WillReturnError(fmt.Errorf("database error"))

		// Call the repository method
		products, totalCount, totalPages, err := productRepo.ShowAllProduct(page, limit)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, totalCount)
		assert.Equal(t, 0, totalPages)
		assert.EqualError(t, err, "database error")
	})
}

func TestGetProductByID(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	productRepo := productrepository.NewProductRepo(db, &log)

	t.Run("Successfully get product by ID", func(t *testing.T) {
		productID := 1

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT $2`)).
			WithArgs(productID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "Product A"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "images" WHERE "images"."product_id" = $1`)).
			WithArgs(productID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).
				AddRow(1, 1).
				AddRow(2, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product_variants" WHERE "product_variants"."product_id" = $1`)).
			WithArgs(productID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product_id"}).
				AddRow(1, 1).
				AddRow(2, 1))

		product, err := productRepo.GetProductByID(productID)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, productID, product.ID)
		assert.Equal(t, "Product A", product.Name)
		assert.Len(t, product.ProductVariant, 2)
		assert.Len(t, product.Image, 2)
	})

	t.Run("Failed to get product by ID due to not found", func(t *testing.T) {
		productID := 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`)).
			WithArgs(productID).
			WillReturnError(fmt.Errorf("record not found"))

		product, err := productRepo.GetProductByID(productID)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "product not found")
	})
}

func TestCreateProduct(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	productRepo := productrepository.NewProductRepo(db, &log)

	t.Run("Successfully create a product", func(t *testing.T) {
		product := &domain.Product{
			Name:        "Product A",
			SKUProduct:  "SKI-2021",
			Price:       100,
			Description: "High-quality product",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Mock database queries
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "products"`).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "product_variants"`).
			WithArgs(
				1,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "images"`).
			WithArgs(
				1,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		err := productRepo.CreateProduct(product)

		assert.NoError(t, err)
		assert.Equal(t, 1, product.ID)
	})

	t.Run("Failed to create product - Product insertion failure", func(t *testing.T) {
		product := &domain.Product{
			Name:        "Product A",
			SKUProduct:  "SKI-2021",
			Price:       100,
			Description: "High-quality product",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "products"`).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("failed to insert product"))

		mock.ExpectRollback()

		err := productRepo.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create product: failed to insert product")
	})

	t.Run("Failed to create product - Product variant insertion failure", func(t *testing.T) {
		product := &domain.Product{
			Name:        "Product A",
			SKUProduct:  "SKI-2021",
			Price:       100,
			Description: "High-quality product",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "products"`).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "product_variants"`).
			WithArgs(
				1,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("failed to insert product variant"))

		mock.ExpectRollback()

		err := productRepo.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create product variant: failed to insert product variant")
	})

	t.Run("Failed to create product - Image insertion failure", func(t *testing.T) {
		product := &domain.Product{
			Name:        "Product A",
			SKUProduct:  "SKI-2021",
			Price:       100,
			Description: "High-quality product",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "products"`).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "product_variants"`).
			WithArgs(
				1,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(`INSERT INTO "images"`).
			WithArgs(
				1,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("failed to insert image"))

		mock.ExpectRollback()

		err := productRepo.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create image: failed to insert image")
	})
}

func TestDeleteProduct(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	productRepo := productrepository.NewProductRepo(db, &log)

	t.Run("Successfully delete product", func(t *testing.T) {
		productID := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "deleted_at"=$1 WHERE "products"."id" = $2 AND "products"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), productID).    // AnyArg untuk timestamp, dan ID produk
			WillReturnResult(sqlmock.NewResult(0, 1)) // Simulasi penghapusan sukses

		mock.ExpectCommit()

		err := productRepo.DeleteProduct(productID)

		assert.NoError(t, err)
	})

	t.Run("Failed to delete product - Product not found", func(t *testing.T) {
		productID := 2

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "deleted_at"=$1 WHERE "products"."id" = $2 AND "products"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), productID).    // `deleted_at` timestamp dan ID produk
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected (no product found)

		mock.ExpectCommit()

		err := productRepo.DeleteProduct(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "product not found")
	})

	t.Run("Failed to delete product - Database error", func(t *testing.T) {
		productID := 3

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "deleted_at"=$1 WHERE "products"."id" = $2 AND "products"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), productID). // `deleted_at` timestamp dan ID produk
			WillReturnError(fmt.Errorf("database error"))

		mock.ExpectRollback()

		err := productRepo.DeleteProduct(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})

}

func TestUpdateProduct(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	productRepo := productrepository.NewProductRepo(db, &log)

	t.Run("Successfully update a product", func(t *testing.T) {
		productID := uint(1)
		product := &domain.Product{
			Name:        "Updated Product",
			SKUProduct:  "SKI-2022",
			Price:       150,
			Description: "Updated description",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "name"=$1,"sku_product"=$2,"price"=$3,"description"=$4,"updated_at"=$5 WHERE id = $6 AND "products"."deleted_at" IS NULL`)).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(), // For updated_at
				productID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected
		mock.ExpectCommit()

		err := productRepo.UpdateProduct(productID, product)

		assert.NoError(t, err)
	})

	t.Run("Failed to update product - No rows affected", func(t *testing.T) {
		productID := uint(2)
		product := &domain.Product{
			Name:        "Another Product",
			SKUProduct:  "SKI-2023",
			Price:       200,
			Description: "Another description",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "name"=$1,"sku_product"=$2,"price"=$3,"description"=$4,"updated_at"=$5 WHERE id = $6 AND "products"."deleted_at" IS NULL`)).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(), // Updated at timestamp
				productID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := productRepo.UpdateProduct(productID, product)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("no record found with shipping_id %d", productID))
	})

	t.Run("Failed to update product - Database error", func(t *testing.T) {
		productID := uint(3)
		product := &domain.Product{
			Name:        "Product with Error",
			SKUProduct:  "SKI-ERROR",
			Price:       300,
			Description: "Error description",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "name"=$1,"sku_product"=$2,"price"=$3,"description"=$4,"updated_at"=$5 WHERE id = $6 AND "products"."deleted_at" IS NULL`)).
			WithArgs(
				product.Name,
				product.SKUProduct,
				product.Price,
				product.Description,
				sqlmock.AnyArg(), // Updated at timestamp
				productID,
			).
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		err := productRepo.UpdateProduct(productID, product)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
