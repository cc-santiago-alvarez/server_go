package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultPort = "3252"

func main() {
	//conexion de base de datoa
	ctx := context.Background()

db, err := mongodb.Connect(ctx, mongodb.Config{
    URI:      os.Getenv("MONGO_URI"),      // la URI con root:12345abc@... vive aqui
    Database: os.Getenv("MONGO_DATABASE"), // tu URI no trae DB, va aparte
})
if err != nil {
    log.Fatal(err) // si Mongo no responde, el proceso no arranca
}
defer db.Close(ctx)

repo := mongocc.NewProductRepository(db.Collection("products"))

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
