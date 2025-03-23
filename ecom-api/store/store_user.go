package store

import (
	"context"
	"fmt"
)

func (s *MySQLStore) CreateUser(ctx context.Context, u *User) (*User, error) {
	query := `
		INSERT INTO users (name, email, password, is_admin)
		VALUES (:name, :email, :password, :is_admin)
	`
	res, err := s.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	u.ID = id

	return u, nil
}

func (s *MySQLStore) GetUser(ctx context.Context, email string) (*User, error) {
	var u User

	if err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email=?", email); err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &u, nil
}

func (s *MySQLStore) ListUsers(ctx context.Context) ([]User, error) {
	var users []User

	if err := s.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	return users, nil
}

func (s *MySQLStore) UpdateUser(ctx context.Context, u *User) (*User, error) {
	query := `UPDATE users SET name=:name, email=:email, password=:password, is_admin=:is_admin, updated_at=:updated_at`
	_, err := s.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

func (s *MySQLStore) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("erorr deleting user: %w", err)
	}

	return nil
}
