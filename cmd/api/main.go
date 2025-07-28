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

	var router http.Handler = router.NewRouter(conn, &cfg.Auth)

	if cfg.Server.Cors {
		router = middleware.CorsMiddleware(router)
	}
	router = middleware.LoggingMiddleware(router)

	fmt.Println("Starting server on port", cfg.Server.Port)
	if cfg.Server.SSL {
		err = http.ListenAndServeTLS(":"+cfg.Server.Port, "./cert.pem", "./key.pem", router)
	} else {
		err = http.ListenAndServe(":"+cfg.Server.Port, router)
	}
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
