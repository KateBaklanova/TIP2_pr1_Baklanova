package main

import (
	"kate/services/auth/internal/http"
	"os"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}
	http.StartServer(port)
}
