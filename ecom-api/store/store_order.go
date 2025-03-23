package store

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (s *MySQLStore) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	err := s.execTx(ctx, func(tx *sqlx.Tx) error {
		// insrt order
		order, err := createOrder(ctx, tx, o)
		if err != nil {
			return fmt.Errorf("error inserting order: %w", err)
		}

		for _, oi := range order.Items {
			oi.OrderID = order.ID
			// insert order items
			if err := createOrderItem(ctx, tx, oi); err != nil {
				return fmt.Errorf("error inserting order items: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	return o, nil
}

func createOrder(ctx context.Context, tx *sqlx.Tx, o *Order) (*Order, error) {
	query := `
		INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, user_id)
		VALUES (:payment_method, :tax_price, :shipping_price, :total_price, :user_id)
	`
	res, err := tx.NamedExecContext(ctx, query, o)
	if err != nil {
		return nil, fmt.Errorf("error inserting order: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", err)
	}
	o.ID = id

	return o, nil
}

func createOrderItem(ctx context.Context, tx *sqlx.Tx, oi OrderItem) error {
	query := `
		INSERT INTO order_items (name, quantity, image, price, product_id, order_id)
		VALUES (:name, :quantity, :image, :price, :product_id, :order_id)
	`
	res, err := tx.NamedExecContext(ctx, query, oi)
	if err != nil {
		return fmt.Errorf("error inserting order items: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	oi.ID = id

	return nil
}

func (s *MySQLStore) GetOrder(ctx context.Context, userID int64) (*Order, error) {
	var o Order
	err := s.db.GetContext(ctx, &o, "SELECT * FROM orders WHERE user_id=?", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}

	var items []OrderItem
	err = s.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", o.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting order items: %w", err)
	}
	o.Items = items

	return &o, nil
}

func (s *MySQLStore) ListOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	if err := s.db.SelectContext(ctx, &orders, "SELECT * FROM orders"); err != nil {
		return nil, fmt.Errorf("error getting orders: %w", err)
	}

	for i := range orders {
		var items []OrderItem

		if err := s.db.SelectContext(ctx, &items, "SELECT * FROM order_items WHERE order_id=?", orders[i].ID); err != nil {
			return nil, fmt.Errorf("error getting order items: %w", err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (s *MySQLStore) UpdateOrderStatus(ctx context.Context, o *Order) (*Order, error) {
	query := "UPDATE orders SET status=:status, updated_at=:updated_at WHERE id=:id"
	_, err := s.db.NamedExecContext(ctx, query, o)
	if err != nil {
		return nil, fmt.Errorf("error updating order status: %w", err)
	}

	return o, nil
}

func (s *MySQLStore) DeleteOrder(ctx context.Context, id int64) error {
	err := s.execTx(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order items: %w", err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM orders WHERE id=?", id)
		if err != nil {
			return fmt.Errorf("error deleting order: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	return nil
}

func (s *MySQLStore) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rollback transaction: %w", err)
		}
		return fmt.Errorf("error in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commit transaction: %w", err)
	}

	return nil
}