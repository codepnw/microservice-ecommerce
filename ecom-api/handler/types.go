package handler

import (
	"time"

	"github.com/codepnw/microservice-ecommerce/ecom-api/server"
	"github.com/codepnw/microservice-ecommerce/token"
)

type handler struct {
	server     *server.Server
	tokenMaker *token.JWTMaker
}

func NewHandler(server *server.Server, secretKey string) *handler {
	return &handler{
		server: server,
		tokenMaker: token.NewJWTMaker(secretKey),
	}
}

type ProductReq struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Rating       int64   `json:"rating"`
	NumReviews   int64   `json:"num_reviews"`
	Price        float64 `json:"price"`
	CountInStock int64   `json:"count_in_stock"`
}

type ProductRes struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Image        string     `json:"image"`
	Category     string     `json:"category"`
	Description  string     `json:"description"`
	Rating       int64      `json:"rating"`
	NumReviews   int64      `json:"num_reviews"`
	Price        float64    `json:"price"`
	CountInStock int64      `json:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

// ========== ORDER ===========
type OrderReq struct {
	ID            int64        `json:"id"`
	Items         []*OrderItem `json:"items"`
	PaymentMethod string       `json:"payment_method"`
	TaxPrice      float32      `json:"tax_price"`
	ShippingPrice float32      `json:"shipping_price"`
	TotalPrice    float32      `json:"total_price"`
	Status        string       `json:"status"`
}

type OrderItem struct {
	Name      string  `json:"name"`
	Quantity  int64   `json:"quantity"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	ProductID int64   `json:"product_id"`
}

type OrderRes struct {
	ID            int64       `json:"id"`
	Items         []OrderItem `json:"items"`
	PaymentMethod string      `json:"payment_method"`
	TaxPrice      float32     `json:"tax_price"`
	ShippingPrice float32     `json:"shipping_price"`
	TotalPrice    float32     `json:"total_price"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     *time.Time  `json:"updated_at"`
}

type UserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserRes struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type ListUserRes struct {
	Users []UserRes `json:"users"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string  `json:"access_token"`
	User        UserRes `json:"user"`
}
