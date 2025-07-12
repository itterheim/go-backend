package core

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) Get(jti string) (*Token, error) {
	token := &Token{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, user_id, jti, created, expiration, blocked
		FROM tokens
		WHERE jti = $1 AND expiration > CURRENT_TIMESTAMP
	`, jti).Scan(&token.ID, &token.UserID, &token.JTI, &token.Created, &token.Expiration, &token.Blocked)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *TokenRepository) Create(token Token) error {
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO tokens (user_id, jti, expiration)
		VALUES ($1, $2, $3)
		RETURNING id
	`, token.UserID, token.JTI, token.Expiration).Scan(&token.ID)

	return err
}

func (r *TokenRepository) Invalidate(id int64) error {
	_, err := r.db.Exec(context.Background(), `
        UPDATE tokens
        SET blocked = TRUE
        WHERE id = $1
    `, id)

	return err
}

func (r *TokenRepository) Delete(jti string) error {
	_, err := r.db.Exec(context.Background(), `
		DELETE FROM tokens
		WHERE jti = $1
	`, jti)

	return err
}
