package handler

import (
	"time"

	"github.com/codepnw/microservice-ecommerce/ecom-api/server"
)

type handler struct {
	server *server.Server
}

func NewHandler(server *server.Server) *handler {
	return &handler{server: server}
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
