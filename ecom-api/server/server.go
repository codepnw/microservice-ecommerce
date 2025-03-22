package server

import (
	"context"

	"github.com/codepnw/microservice-ecommerce/ecom-api/store"
)

type Server struct {
	store *store.MySQLStore
}

func NewServer(store *store.MySQLStore) *Server {
	return &Server{store: store}
}

// ========= PRODUCT ==========
func (s *Server) CreateProduct(ctx context.Context, p *store.Product) (*store.Product, error) {
	return s.store.CreateProduct(ctx, p)
}

func (s *Server) GetProduct(ctx context.Context, id int64) (*store.Product, error) {
	return s.store.GetProduct(ctx, id)
}

func (s *Server) ListProducts(ctx context.Context) ([]store.Product, error) {
	return s.store.ListProducts(ctx)
}

func (s *Server) UpdateProduct(ctx context.Context, p *store.Product) (*store.Product, error) {
	return s.store.UpdateProduct(ctx, p)
}

func (s *Server) DeleteProduct(ctx context.Context, id int64) error {
	return s.store.DeleteProduct(ctx, id)
}

// ========= ORDER ==========
func (s *Server) CreateOrder(ctx context.Context, o *store.Order) (*store.Order, error) {
	return s.store.CreateOrder(ctx, o)
}

func (s *Server) GetOrder(ctx context.Context, id int64) (*store.Order, error) {
	return s.store.GetOrder(ctx, id)
}

func (s *Server) ListOrder(ctx context.Context) ([]store.Order, error) {
	return s.store.ListOrders(ctx)
}

func (s *Server) DeleteOrder(ctx context.Context, id int64) error {
	return s.store.DeleteOrder(ctx, id)
}