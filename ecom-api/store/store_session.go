package store

import (
	"context"
	"fmt"
)

func (s *MySQLStore) CreateSession(ctx context.Context, sess *Session) (*Session, error) {
	query := `
		INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at)
		VALUES (:id, :user_email, :refresh_token, :is_revoked, :expires_at)	
	`
	_, err := s.db.NamedExecContext(ctx, query, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *MySQLStore) GetSession(ctx context.Context, id string) (*Session, error) {
	var session Session
	err := s.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &session, nil
}

func (s *MySQLStore) RevokeSession(ctx context.Context, id string) error {
	query := "UPDATE sessions SET is_revoked=1 WHERE id=:id"
	_, err := s.db.NamedExecContext(ctx, query, map[string]any{"id": id})
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) DeleteSession(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
