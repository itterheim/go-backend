package core

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActionRepository struct {
	db *pgxpool.Pool
}

func NewActionRepository(db *pgxpool.Pool) *ActionRepository {

	return &ActionRepository{db: db}
}

func (r *ActionRepository) ListActions() ([]Action, error) {
	return nil, errors.New("Not implemented")
}

func (r *ActionRepository) GetAction(id int64) (*Action, error) {
	// tx, err := db.BeginTx(context.Background(), pgx.TxOptions{})

	return nil, errors.New("Not implemented")
}

func (r *ActionReference) CreateAction(data *Action) (*Action, error) {
	return nil, errors.New("Not implemented")
}

func (r *ActionRepository) UpdateAction(data *Action) (*Action, error) {
	return nil, errors.New("Not implemented")
}

func (r *ActionRepository) DeleteAction(id int64) error {
	return errors.New("Not implemented")
}
