package core

import (
	"backend/pkg/jwt"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(id int64) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, created, updated, username, password, role
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Password, &user.Role)
	if err == pgx.ErrNoRows {
		fmt.Println(err)
		return nil, err
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, created, updated, username, password, role
		FROM users
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Password, &user.Role)
	if err == pgx.ErrNoRows {
		fmt.Println(err)
		return nil, err
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(username, password string, role jwt.ClaimRole) (*User, error) {
	user := &User{}

	err := r.db.QueryRow(context.Background(), `
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, created, updated, username, password, role
	`, username, password, role).Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *UserRepository) ListUsers() ([]User, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT
			id, created, updated, username, role
		FROM users
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
