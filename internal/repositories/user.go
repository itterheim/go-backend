package repositories

import (
	"backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *User {
	return &User{db: db}
}

func (r *User) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id, created, updated, username, password
		FROM users
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Password)
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

func (r *User) Create(username, password string) (*models.User, error) {
	user := &models.User{}

	err := r.db.QueryRow(context.Background(), `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id, created, updated, username, password
	`, username, password).Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *User) ListUsers() ([]models.User, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT
			id, created, updated, username, password
		FROM users
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.Created, &user.Updated, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
