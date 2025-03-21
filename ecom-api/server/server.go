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