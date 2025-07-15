package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/router"
	"backend/pkg/middleware"
	"fmt"
	"net/http"
	"time"
)

func main() {
	startTime := time.Now()
	defer func() {
		endTime := time.Since(startTime)
		fmt.Println("Done in:", endTime.Seconds(), "s")
	}()

	cfg, err := config.LoadConfig("./config/")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	conn, err := db.ConnectPgx(cfg)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer conn.Close()

	err = db.Check(conn)
	if err != nil {
		panic(err)
	}

	mux := router.NewRouter(conn, &cfg.Auth)

	router := middleware.LoggingMiddleware(mux)

	fmt.Println("Starting server on port", cfg.Server.Port)
	err = http.ListenAndServe(":"+cfg.Server.Port, router)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
