package db

import (
	"backend/internal/config"
	"context"
	"errors"
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

func Check(db *pgxpool.Pool) error {
	rows, err := db.Query(context.Background(), `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema='public' AND table_type='BASE TABLE'
	`)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Failed to list databse tables")
	}

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Println("Error scanning table name:", err.Error())
			return errors.New("Failed to scan table name")
		}
		fmt.Println("Table:", tableName)
	}

	return nil
}
