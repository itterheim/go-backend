package middleware

import "net/http"

// corsMiddleware is a simple middleware to handle CORS.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set common CORS headers.
		// NOTE: In a production environment, you should replace "*" with specific origins.
		// For example: "http://localhost:3000"

		w.Header().Set("Access-Control-Allow-Origin", "https://localhost:5173") // "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Set to "true" if you need to allow cookies, etc.
		w.Header().Set("Access-Control-Max-Age", "300")            // Cache preflight results for 5 minutes

		// Handle preflight OPTIONS requests.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
