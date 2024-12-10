package database

import (
	"project/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	var err error

	if err = dropTables(db); err != nil {
		return err
	}

	if err = setupJoinTables(db); err != nil {
		return err
	}

	if err = autoMigrates(db); err != nil {
		return err
	}

	return createViews(db)
}

func autoMigrates(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.PasswordResetToken{},
		&domain.Product{},
		&domain.ProductVariant{},
		&domain.Image{},
		&domain.Customer{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Review{},
		&domain.Stock{},
		&domain.Promotion{},
	)
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&domain.User{},
		&domain.Category{},
		&domain.PasswordResetToken{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Customer{},
		&domain.Product{},
		&domain.ProductVariant{},
		&domain.Image{},
		&domain.Review{},
		&domain.Stock{},
		&domain.Promotion{},
	)
}

func setupJoinTables(db *gorm.DB) error {
	var err error

	return err
}

func createViews(db *gorm.DB) error {
	var err error
	if err = queryOrders(db); err != nil {
		return err
	}

	if err = queryOrderItems(db); err != nil {
		return err
	}

	return err
}

func queryOrders(db *gorm.DB) error {
	query := db.Raw(`
		SELECT orders.id, customers.name AS customer_name, customers.address AS customer_address, orders.payment_method, orderitems.total, orders.status
		FROM orders
		JOIN (
			SELECT order_id, SUM(quantity * unit_price) AS total
			FROM order_items
			GROUP BY order_id) orderitems ON orders.id = orderitems.order_id
		JOIN customers ON orders.customer_id = customers.id
	`)
	return db.Migrator().CreateView("order_totals", gorm.ViewOption{Query: query, Replace: true})
}

func queryOrderItems(db *gorm.DB) error {
	query := db.Raw(`
		SELECT order_items.order_id, CONCAT_WS(' ',products.name, 'size: '||product_variants.size, 'color: '||product_variants.color) AS product_name, order_items.quantity, order_items.unit_price, order_items.quantity * order_items.unit_price AS subtotal
		FROM order_items
		JOIN product_variants ON order_items.variant_id = product_variants.id
		JOIN products ON product_variants.product_id = products.id
	`)
	return db.Migrator().CreateView("order_item_subtotals", gorm.ViewOption{Query: query, Replace: true})
}
