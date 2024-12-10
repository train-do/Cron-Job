package dashboardrepository_test

import (
	"regexp"
	"testing"

	"project/helper"
	dashboardrepository "project/repository/dashboard_repository"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestGetEarningDashboard(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := zap.NewNop()

	dashboardRepo := dashboardrepository.NewDashboardRepo(db, log)

	rows := sqlmock.NewRows([]string{"total_amount"}).AddRow(500000)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT SUM(oi.unit_price * oi.quantity) as total_amount FROM orders as o 
		JOIN order_items as oi ON oi.order_id = o.id 
		WHERE o.status = $1 AND (o.created_at BETWEEN $2 AND $3)`)).
		WithArgs("completed", sqlmock.AnyArg(), sqlmock.AnyArg()). // Use sqlmock.AnyArg() to match any argument for time
		WillReturnRows(rows)

	totalEarning, err := dashboardRepo.GetEarningDashboard()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if totalEarning != 500000 {
		t.Errorf("Expected total earning 500000, but got %d", totalEarning)
	}

	if err != nil {
		t.Fatalf("Error should be nil, but got: %v", err)
	}

	if totalEarning != 500000 {
		t.Fatalf("Expected total earning 500000, but got: %d", totalEarning)
	}
}

func TestGetSummary(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := zap.NewNop()

	dashboardRepo := dashboardrepository.NewDashboardRepo(db, log)

	rows := sqlmock.NewRows([]string{"sales", "orders", "items"}).AddRow(1000000, 50, 150)

	mockUserRows := sqlmock.NewRows([]string{"users"}).AddRow(30)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT SUM(oi.unit_price * oi.quantity) as sales, COUNT(o.id) as orders, SUM(oi.quantity) as items FROM orders as o 
		JOIN order_items as oi ON oi.order_id = o.id 
		WHERE o.status != $1 AND (o.created_at BETWEEN $2 AND $3)`)).
		WithArgs("canceled", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT COUNT(DISTINCT o.customer_id) as users FROM orders as o 
		WHERE o.created_at BETWEEN $1 AND $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mockUserRows)

	summary, err := dashboardRepo.GetSummary()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if summary.Sales != 1000000 {
		t.Errorf("Expected sales 1000000, but got %d", summary.Sales)
	}
	if summary.Orders != 50 {
		t.Errorf("Expected 50 orders, but got %d", summary.Orders)
	}
	if summary.Items != 150 {
		t.Errorf("Expected 150 items, but got %d", summary.Items)
	}
	if summary.Users != 30 {
		t.Errorf("Expected 30 users, but got %d", summary.Users)
	}
}

func TestGetBestSeller(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := zap.NewNop()
	dashboardRepo := dashboardrepository.NewDashboardRepo(db, log)

	rows := sqlmock.NewRows([]string{"product_id", "total_sold"}).
		AddRow(1, 100).
		AddRow(2, 80).
		AddRow(3, 60).
		AddRow(4, 50).
		AddRow(5, 40)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT oi.variant_id as product_id, 
		SUM(oi.quantity) as total_sold 
		FROM orders as o JOIN order_items as oi ON oi.order_id = o.id 
		WHERE o.status = $1 AND (o.created_at BETWEEN $2 AND $3) 
		GROUP BY "oi"."variant_id" ORDER BY total_sold DESC LIMIT $4`)).
		WithArgs("completed", sqlmock.AnyArg(), sqlmock.AnyArg(), 5).
		WillReturnRows(rows)

	bestSellers, err := dashboardRepo.GetBestSeller()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(bestSellers) != 5 {
		t.Errorf("Expected 5 best sellers, but got %d", len(bestSellers))
	}

	if bestSellers[0].ProductID != 1 || bestSellers[0].TotalSold != 100 {
		t.Errorf("Expected best seller product ID 1 with total sold 100, but got product ID %d with total sold %d", bestSellers[0].ProductID, bestSellers[0].TotalSold)
	}

	if bestSellers[1].ProductID != 2 || bestSellers[1].TotalSold != 80 {
		t.Errorf("Expected second best seller product ID 2 with total sold 80, but got product ID %d with total sold %d", bestSellers[1].ProductID, bestSellers[1].TotalSold)
	}
}

func TestGetMonthlyRevenue(t *testing.T) {
	db, mock := helper.SetupTestDB()
	defer func() { _ = mock.ExpectationsWereMet() }()

	log := zap.NewNop()
	dashboardRepo := dashboardrepository.NewDashboardRepo(db, log)

	rows := sqlmock.NewRows([]string{"month", "revenue"}).
		AddRow("January", 500000).
		AddRow("February", 600000).
		AddRow("March", 700000)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT TO_CHAR(o.created_at, 'Month') as month, SUM(oi.unit_price * oi.quantity) as revenue FROM orders as o 
		JOIN order_items as oi ON oi.order_id = o.id 
		WHERE o.status = $1 AND o.created_at >= $2 
		GROUP BY TO_CHAR(o.created_at, 'Month') ORDER BY MIN(o.created_at)`)).
		WithArgs("completed", sqlmock.AnyArg()).
		WillReturnRows(rows)

	revenues, err := dashboardRepo.GetMonthlyRevenue()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(revenues) != 3 {
		t.Errorf("Expected 3 monthly revenues, but got %d", len(revenues))
	}

	if revenues[0].Month != "January" || revenues[0].Revenue != 500000 {
		t.Errorf("Expected January revenue 500000, but got month %s with revenue %d", revenues[0].Month, revenues[0].Revenue)
	}

	if revenues[1].Month != "February" || revenues[1].Revenue != 600000 {
		t.Errorf("Expected February revenue 600000, but got month %s with revenue %d", revenues[1].Month, revenues[1].Revenue)
	}
}
