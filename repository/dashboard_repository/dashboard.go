package dashboardrepository

import (
	"project/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardRepo interface {
	GetEarningDashboard() (int, error)
	GetSummary() (*domain.Summary, error)
	GetBestSeller() ([]*domain.BestSeller, error)
	GetMonthlyRevenue() ([]*domain.Revenue, error)
}

type dashboardRepo struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewDashboardRepo(db *gorm.DB, log *zap.Logger) DashboardRepo {

	return &dashboardRepo{db, log}
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, -1).Add(24*time.Hour - time.Second)
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

func (dr *dashboardRepo) GetEarningDashboard() (int, error) {

	var totalEarning int
	now := time.Now()

	query := dr.db.Table("orders as o").
		Select("SUM(oi.unit_price * oi.quantity) as total_amount").
		Joins("JOIN order_items as oi ON oi.order_id = o.id").
		Where("o.status = ?", "completed").
		Where("o.created_at BETWEEN ? AND ?",
			StartOfMonth(now),
			EndOfMonth(now),
		).
		Scan(&totalEarning)

	if query.Error != nil {
		return 0, query.Error
	}

	return totalEarning, nil
}

func (dr *dashboardRepo) GetSummary() (*domain.Summary, error) {
	var summary domain.Summary
	now := time.Now()

	query := dr.db.Table("orders as o").
		Select("SUM(oi.unit_price * oi.quantity) as sales, COUNT(o.id) as orders, SUM(oi.quantity) as items").
		Joins("JOIN order_items as oi ON oi.order_id = o.id").
		Where("o.status != ?", "canceled").
		Where("o.created_at BETWEEN ? AND ?", StartOfMonth(now),
			EndOfMonth(now)).
		Scan(&summary)

	if query.Error != nil {
		return nil, query.Error
	}

	var users int
	queryUser := dr.db.Table("orders as o").
		Select("COUNT(DISTINCT o.customer_id) as users").
		Where("o.created_at BETWEEN ? AND ?", StartOfMonth(now),
			EndOfMonth(now)).
		Scan(&users)

	if queryUser.Error != nil {
		return nil, queryUser.Error
	}

	summary.Users = users

	return &summary, nil
}

func (dr *dashboardRepo) GetBestSeller() ([]*domain.BestSeller, error) {
	var bestSellers []*domain.BestSeller
	now := time.Now()

	query := dr.db.Table("orders as o").
		Select("oi.variant_id as product_id, SUM(oi.quantity) as total_sold").
		Joins("JOIN order_items as oi ON oi.order_id = o.id").
		Group("oi.variant_id").
		Where("o.status = ?", "completed").
		Where("o.created_at BETWEEN ? AND ?", StartOfMonth(now), EndOfMonth(now)).
		Order("total_sold DESC").
		Limit(5).
		Scan(&bestSellers)

	if query.Error != nil {
		return nil, query.Error
	}

	return bestSellers, nil
}

func (dr *dashboardRepo) GetMonthlyRevenue() ([]*domain.Revenue, error) {
	var revenues []*domain.Revenue
	now := time.Now()

	query := dr.db.Table("orders as o").
		Select("TO_CHAR(o.created_at, 'Month') as month, SUM(oi.unit_price * oi.quantity) as revenue").
		Joins("JOIN order_items as oi ON oi.order_id = o.id").
		Where("o.status = ?", "completed").
		Where("o.created_at >= ?", StartOfYear(now)).
		Group("TO_CHAR(o.created_at, 'Month')").
		Order("MIN(o.created_at)").
		Scan(&revenues)

	if query.Error != nil {
		return nil, query.Error
	}

	return revenues, nil
}
