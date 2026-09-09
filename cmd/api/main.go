package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = "3252"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	addr := ":" + port

	fmt.Println("Server running on", addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

	log.Fatal(http.ListenAndServe(addr, mux))
}

func health(w http.ResponseWriter, _ *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("everything is ok!"))
}
