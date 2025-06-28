package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/cli"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var configPath string
var migrationsPath string

func main() {
	startTime := time.Now()
	defer func() {
		endTime := time.Since(startTime)
		fmt.Println("Done in:", endTime.Seconds(), "s")
	}()

	configPath, err := cli.GetArg("config")
	if err != nil {
		configPath = "../../config/"
	}

	migrationsPath, err = cli.GetArg("migrations")
	if err != nil {
		migrationsPath = "../../migrations/"
	}

	// load and parse the configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Println(cfg)

	// connect to the database
	conn, err := db.ConnectPgx(cfg)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer conn.Close()

	// Clean up the database
	_, err = cli.GetArg("nuke")
	if err == nil {
		err = nukeDatabase(conn)
		if err != nil {
			fmt.Println("Error nuking database:", err)
			return
		}
	}

	// read version from the database
	version, err := getSchemaVersion(conn)
	if err != nil {
		fmt.Println("Error reading schema version:", err)
		initDatabase(conn)
		version, err = getSchemaVersion(conn)
		if err != nil {
			fmt.Println("Error reading schema version after init:", err)
			return
		}
	}

	fmt.Println("Current schema version:", version)

	// migrate the database
	err = migrate(conn, version)
	if err != nil {
		fmt.Println("Error migrating database:", err)
		return
	}

	fmt.Println("Database migrated successfully")

	// Setup the initial user
	username, err := cli.GetArg("username")
	if err != nil || len(username) == 0 {
		fmt.Println("no default user defined")
		return
	}

	password, err := cli.GetArg("password")
	if err != nil {
		fmt.Println("no default user password defined")
	}

	err = createUser(conn, username, password, cfg)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}
}

func nukeDatabase(conn *pgxpool.Pool) error {
	// load the nuke sql file
	nukeSQL, err := loadFile(migrationsPath + "nuke.sql")
	if err != nil {
		return err
	}

	// execute the nuke sql
	_, err = conn.Exec(context.Background(), nukeSQL)
	if err != nil {
		return err
	}

	fmt.Println("Database nuked successfully")

	return nil
}

func getSchemaVersion(conn *pgxpool.Pool) (string, error) {
	var version string
	err := conn.QueryRow(context.Background(), `
		SELECT version
		FROM migrations
		ORDER BY timestamp DESC
		LIMIT 1
	`).Scan(&version)

	if err == pgx.ErrNoRows {
		const initVersion = "00000000-00"
		// insert 000 version
		_, err = conn.Exec(context.Background(), `
			INSERT INTO migrations (version)
			VALUES ('`+initVersion+`')
		`)
		if err != nil {
			return "", err
		}
		return initVersion, nil
	}

	if err != nil {
		return "", err
	}

	return version, nil
}

func initDatabase(conn *pgxpool.Pool) error {
	initSQL, err := loadFile(migrationsPath + "init.sql")
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), initSQL)
	if err != nil {
		return err
	}

	return nil
}

func migrate(conn *pgxpool.Pool, version string) error {
	files, err := getMigrationFiles(version)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = runMigration(conn, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func runMigration(conn *pgxpool.Pool, file string) error {
	// load migration file
	migration, err := os.ReadFile(migrationsPath + file)
	if err != nil {
		fmt.Println("Error loading migration file:", err)
		return err
	}

	// execute migration
	fmt.Println("Executing migration file:", file)

	_, err = conn.Exec(context.Background(), string(migration))
	if err != nil {
		fmt.Println("Error executing migration file:", err)
		return err
	}

	version, err := parseVersion(file)
	if err != nil {
		fmt.Println("Error parsing version from file:", err)
		return err
	}

	// insert version
	_, err = conn.Exec(context.Background(), `
		INSERT INTO migrations (version)
		VALUES ($1)
	`, version)
	if err != nil {
		fmt.Println("Error inserting version", file[:3], ":", err)
		return err
	}

	return nil
}

func getMigrationFiles(version string) ([]string, error) {
	// list all migration files
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		fmt.Println("Error reading migrations directory:", err)
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		name := file.Name()
		fileVersion, err := parseVersion(name)
		if err != nil {
			continue
		}

		if fileVersion > version {
			migrationFiles = append(migrationFiles, name)
		}
	}

	sort.Strings(migrationFiles)

	return migrationFiles, nil
}

func loadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func parseVersion(filename string) (string, error) {
	if !strings.Contains(filename, ".sql") {
		return "", fmt.Errorf("invalid filename %s", filename)
	}

	parts := strings.Split(filename[:len(filename)-4], "-")
	if len(parts) >= 2 {
		return parts[0] + "-" + parts[1], nil
	}
	return "", fmt.Errorf("invalid filename %s", filename)
}

func createUser(conn *pgxpool.Pool, username, password string, config *config.Config) error {
	userRepo := repositories.NewUserRepository(conn)
	user, err := userRepo.GetByUsername(username)
	if user != nil {
		fmt.Println("User already exists")
		return nil
	}

	user, err = models.NewUser(username, password, config.Auth.BCryptCost)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return err
	}

	user, err = userRepo.Create(user.Username, user.Password)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return err
	}

	fmt.Println("User created:", user.ID)
	return nil
}
