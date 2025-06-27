package repositories

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Token struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *Token {
	return &Token{
		db: db,
	}
}

func (r *Token) Get(jti string) (*models.Token, error) {
	token := &models.Token{}
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

func (r *Token) Create(token models.Token) error {
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO tokens (user_id, jti, expiration)
		VALUES ($1, $2, $3)
		RETURNING id
	`, token.UserID, token.JTI, token.Expiration).Scan(&token.ID)

	return err
}

func (r *Token) Invalidate(id int64) error {
	_, err := r.db.Exec(context.Background(), `
        UPDATE tokens
        SET blocked = TRUE
        WHERE id = $1
    `, id)

	return err
}

func (r *Token) Delete(jti string) error {
	_, err := r.db.Exec(context.Background(), `
		DELETE FROM tokens
		WHERE jti = $1
	`, jti)

	return err
}
