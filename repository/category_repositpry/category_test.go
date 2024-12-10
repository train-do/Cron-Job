package categoryrepositpry_test

import (
	"fmt"
	"project/domain"
	"project/helper"
	categoryrepositpry "project/repository/category_repositpry"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestShowAllCategory(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	categoryRepo := categoryrepositpry.NewCategoryRepo(db, &log)

	t.Run("Successfully show all categories", func(t *testing.T) {
		page, limit := 1, 2

		// Mock count query
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" LIMIT $1`)).
			WithArgs(limit).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "Category 1").
				AddRow(2, "Category 2"))

		categories, totalCount, totalPages, err := categoryRepo.ShowAllCategory(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		assert.Len(t, *categories, 2)
		assert.Equal(t, "Category 1", (*categories)[0].Name)
		assert.Equal(t, "Category 2", (*categories)[1].Name)
		assert.Equal(t, 5, totalCount)
		assert.Equal(t, 3, totalPages)
	})

	t.Run("Category not found", func(t *testing.T) {
		page, limit := 1, 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" LIMIT $1`)).
			WithArgs(limit).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		categories, totalCount, totalPages, err := categoryRepo.ShowAllCategory(page, limit)

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Equal(t, 0, totalCount)
		assert.Equal(t, 0, totalPages)
		assert.EqualError(t, err, "category not found")
	})

	t.Run("Failed to show all categories due to database error", func(t *testing.T) {
		page, limit := 1, 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories"`)).
			WillReturnError(fmt.Errorf("database error"))

		categories, totalCount, totalPages, err := categoryRepo.ShowAllCategory(page, limit)

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Equal(t, 0, totalCount)
		assert.Equal(t, 0, totalPages)
		assert.EqualError(t, err, "database error")
	})

	t.Run("Pagination behavior with multiple pages", func(t *testing.T) {
		page, limit := 2, 2

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "categories"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" LIMIT $1 OFFSET $2`)).
			WithArgs(limit, (page-1)*limit).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(3, "Category 3").
				AddRow(4, "Category 4"))

		categories, totalCount, totalPages, err := categoryRepo.ShowAllCategory(page, limit)

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		assert.Len(t, *categories, 2)
		assert.Equal(t, "Category 3", (*categories)[0].Name)
		assert.Equal(t, "Category 4", (*categories)[1].Name)
		assert.Equal(t, 5, totalCount)
		assert.Equal(t, 3, totalPages)
	})
}

func TestDeleteCategory(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	categoryRepo := categoryrepositpry.NewCategoryRepo(db, &log)

	t.Run("Successfully delete category", func(t *testing.T) {
		id := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := categoryRepo.DeleteCategory(id)

		assert.NoError(t, err)
	})

	t.Run("Category not found", func(t *testing.T) {
		id := 999

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := categoryRepo.DeleteCategory(id)

		assert.Error(t, err)
		assert.EqualError(t, err, "category not found")
	})
}

func TestCreateCategory(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	categoryRepo := categoryrepositpry.NewCategoryRepo(db, &log)

	t.Run("Successfully create category", func(t *testing.T) {
		category := &domain.Category{
			Name: "New Category",
			Icon: "http//skall.jpg",
		}

		mock.ExpectBegin()
		mock.ExpectQuery((`INSERT INTO "categories"`)).
			WithArgs(
				category.Name,
				category.Icon,
				sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := categoryRepo.CreateCategory(category)

		// Assertions
		assert.NoError(t, err)
	})

	t.Run("Failed to create category due to database error", func(t *testing.T) {
		category := &domain.Category{
			Name: "New Category",
			Icon: "http//skall.jpg",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "categories"`).
			WithArgs(
				category.Name,
				category.Icon,
				sqlmock.AnyArg(),
				sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		err := categoryRepo.CreateCategory(category)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create category: database error")
	})
}

func TestGetCategoryByID(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	categoryRepo := categoryrepositpry.NewCategoryRepo(db, &log)

	t.Run("Successfully get category by ID", func(t *testing.T) {
		category := &domain.Category{
			ID:   1,
			Name: "Category 1",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(category.ID, category.Name))

		result, err := categoryRepo.GetCategoryByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, category.ID, result.ID)
		assert.Equal(t, category.Name, result.Name)
	})

	t.Run("Category not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		result, err := categoryRepo.GetCategoryByID(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "category not found or already deleted")
	})

	t.Run("Failed to get category due to database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		result, err := categoryRepo.GetCategoryByID(1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "database error")
	})
}

func TestUpdateCategory(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := *zap.NewNop()
	categoryRepo := categoryrepositpry.NewCategoryRepo(db, &log)

	t.Run("Successfully update category", func(t *testing.T) {
		category := &domain.Category{
			ID:   1,
			Name: "Updated Category",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "id" = $4`)).
			WithArgs(
				category.Name,
				sqlmock.AnyArg(),
				category.ID,
				category.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := categoryRepo.UpdateCategory(int(category.ID), category)

		assert.NoError(t, err)
	})

	t.Run("Category not found", func(t *testing.T) {
		category := &domain.Category{
			ID:   1,
			Name: "Updated Category",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "id" = $4`)).
			WithArgs(
				category.Name,
				sqlmock.AnyArg(),
				category.ID,
				category.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := categoryRepo.UpdateCategory(int(category.ID), category)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("no record found with shipping_id %d", category.ID))
	})

	t.Run("Failed to update category due to database error", func(t *testing.T) {
		category := &domain.Category{
			ID:   1,
			Name: "Updated Category",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "name"=$1,"updated_at"=$2 WHERE id = $3 AND "id" = $4`)).
			WithArgs(
				category.Name,
				sqlmock.AnyArg(),
				category.ID,
				category.ID,
			).
			WillReturnError(fmt.Errorf("database error"))
		mock.ExpectRollback()

		err := categoryRepo.UpdateCategory(int(category.ID), category)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
