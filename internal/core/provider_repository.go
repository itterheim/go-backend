package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProviderRepository struct {
	db *pgxpool.Pool
}

func NewProviderRepository(db *pgxpool.Pool) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) GetById(id int64) (*Provider, error) {
	provider := &Provider{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, created, updated,  name, description, expiration
		FROM providers
		WHERE id = $1
	`, id).Scan(
		&provider.ID,
		&provider.Created,
		&provider.Updated,
		&provider.Name,
		&provider.Description,
		&provider.Expiration,
	)
	if err == pgx.ErrNoRows {
		fmt.Println("Provider repository: provider not found", err)
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (r *ProviderRepository) Create(userID int64, name, description string) (*Provider, error) {
	provider := &Provider{}

	err := r.db.QueryRow(context.Background(), `
		INSERT INTO providers (name, description)
		VALUES ($1, $2)
		RETURNING id, created, updated, name, description
	`, name, description).Scan(
		&provider.ID,
		&provider.Created,
		&provider.Updated,
		&provider.Name,
		&provider.Description,
	)
	if err != nil {
		return nil, err
	}

	return provider, err
}

func (r *ProviderRepository) Update(data *Provider) (*Provider, error) {
	provider := &Provider{}

	err := r.db.QueryRow(context.Background(), `
		UPDATE providers
			SET name = $2, description = $3
		WHERE id = $1
		RETURNING id, created, updated, name, description, expiration
	`, data.ID, data.Name, data.Description).Scan(
		&provider.ID,
		&provider.Created,
		&provider.Updated,
		&provider.Name,
		&provider.Description,
		&provider.Expiration,
	)
	if err != nil {
		return nil, err
	}

	return provider, err
}

func (r *ProviderRepository) Delete(providerId int64) error {
	commandTag, err := r.db.Exec(context.Background(), `
		DELETE FROM providers
		WHERE id = $1
	`, providerId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("Delete: no rows updated")
	}

	return nil
}

func (r *ProviderRepository) List() ([]Provider, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT
			id, created, updated, name, description, expiration
		FROM providers
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := make([]Provider, 0)
	for rows.Next() {
		provider := Provider{}
		err = rows.Scan(
			&provider.ID,
			&provider.Created,
			&provider.Updated,
			&provider.Name,
			&provider.Description,
			&provider.Expiration,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

func (r *ProviderRepository) UpdateToken(providerId int64, jti string, expiration time.Time) error {
	commandTag, err := r.db.Exec(context.Background(), `
		UPDATE providers
		SET jti = $2, expiration = $3
		WHERE id = $1
	`, providerId, jti, expiration)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("UpdateToken: no rows updated")
	}
	return nil
}

func (r *ProviderRepository) RevokeToken(providerId int64) error {
	commandTag, err := r.db.Exec(context.Background(), `
		UPDATE providers
		SET jti = NULL, expiration = NULL
		WHERE id = $1
	`, providerId)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("RevokeToken: no rows updated")
	}
	return nil
}
