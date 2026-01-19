package config

import (
	"fmt"
	"net/http"
)

func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			fmt.Println(origin)
			if origin == allowedOrigin {
				fmt.Println("allowed Origins")
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Expose-Headers", "Location")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				fmt.Println("Options")
				w.WriteHeader(http.StatusNoContent)
				// return
			}

			next.ServeHTTP(w, r)
		})
	}
}
