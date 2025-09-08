package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MetaDandy/Assistense-System/config"
	"github.com/MetaDandy/Assistense-System/src"
)

func main() {
	config.Load()

	router := src.SetupRoutes()

	handler := corsMiddleware(router)

	port := config.Port
	fmt.Printf("Servidor corriendo en puerto %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
