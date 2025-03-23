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

// ========= USER ==========
func (s *Server) CreateUser(ctx context.Context, u *store.User) (*store.User, error) {
	return s.store.CreateUser(ctx, u)
}

func (s *Server) GetUser(ctx context.Context, email string) (*store.User, error) {
	return s.store.GetUser(ctx, email)
}

func (s *Server) ListUsers(ctx context.Context) ([]store.User, error) {
	return s.store.ListUsers(ctx)
}

func (s *Server) UpdateUser(ctx context.Context, u *store.User) (*store.User, error) {
	return s.store.UpdateUser(ctx, u)
}

func (s *Server) DeleteUser(ctx context.Context, id int64) error {
	return s.store.DeleteUser(ctx, id)
}

// ========= SESSION ==========
func (s *Server) CreateSession(ctx context.Context, sess *store.Session) (*store.Session, error) {
	return s.store.CreateSession(ctx, sess)
}

func (s *Server) GetSession(ctx context.Context, id string) (*store.Session, error) {
	return s.store.GetSession(ctx, id)
}

func (s *Server) RevokeSession(ctx context.Context, id string) error {
	return s.store.RevokeSession(ctx, id)
}

func (s *Server) DeleteSession(ctx context.Context, id string) error {
	return s.store.DeleteSession(ctx, id)
}
