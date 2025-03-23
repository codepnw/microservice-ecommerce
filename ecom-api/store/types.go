package store

import "time"

type Product struct {
	ID           int64      `db:"id"`
	Name         string     `db:"name"`
	Image        string     `db:"image"`
	Category     string     `db:"category"`
	Description  string     `db:"description"`
	Rating       int64      `db:"rating"`
	NumReviews   int64      `db:"num_reviews"`
	Price        float64    `db:"price"`
	CountInStock int64      `db:"count_in_stock"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type Order struct {
	ID            int64      `db:"id"`
	PaymentMethod string     `db:"payment_method"`
	TaxPrice      float32    `db:"tax_price"`
	ShippingPrice float32    `db:"shipping_price"`
	TotalPrice    float32    `db:"total_price"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	Items         []OrderItem
}

type OrderItem struct {
	ID        int64   `db:"id"`
	Name      string  `db:"name"`
	Quantity  int64   `db:"quantity"`
	Image     string  `db:"image"`
	Price     float64 `db:"price"`
	ProductID int64   `db:"product_id"`
	OrderID   int64   `db:"order_id"`
}

type User struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updatedAt"`
}
