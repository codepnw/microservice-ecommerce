package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MySQLStore struct {
	db *sqlx.DB
}

func NewMySQLStore(db *sqlx.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

func (s *MySQLStore) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	query := `
		INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) 
		VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)
	`
	res, err := s.db.NamedExecContext(ctx, query, p)
	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", err)
	}
	p.ID = id

	return p, nil
}

func (s *MySQLStore) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product
	query := `SELECT * FROM products WHERE id=?`
	if err := s.db.GetContext(ctx, &p, query, id); err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	return &p, nil
}

func (s *MySQLStore) ListProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	query := `SELECT * FROM products`
	if err := s.db.SelectContext(ctx, &products, query); err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return products, nil
}

func (s *MySQLStore) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	query := `
		UPDATE products 
		SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock
		WHERE id=:id
	`
	if _, err := s.db.NamedExecContext(ctx, query, p); err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return p, nil
}

func (s *MySQLStore) DeleteProduct(ctx context.Context, id int64) error {	
	query := `DELETE FROM products WHERE id=?`
	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}