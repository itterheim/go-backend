package db

import (
	"backend/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPgx(config *config.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Database.GetConnectionString())
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return nil, err
	}
	fmt.Println("Connected to the database successfully")

	return conn, err
}
