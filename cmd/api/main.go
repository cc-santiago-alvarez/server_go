package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"server_go/internal/config"
	"server_go/internal/core/demo/application"
	httpDir "server_go/internal/core/demo/infrastructure/http"
	demomongo "server_go/internal/core/demo/infrastructure/mongo"
	mongodb "server_go/internal/database/mongo"
)

const defaultPort = "3252"

func main() {
	//Esto carga las variables de entorno
	cfg := config.InitConfig()
	cfg.LoadEnvs()

	//conexion de base de datos
	ctx := context.Background()

	db, err := mongodb.Connect(ctx, mongodb.Config{
		URI:      os.Getenv("MONGO_URI"),      // la URI con root:12345abc@... vive aqui
		Database: os.Getenv("MONGO_DATABASE"), // tu URI no trae DB, va aparte
	})
	if err != nil {
		log.Fatal(err) // si Mongo no responde, el proceso no arranca
	}
	defer db.Close(ctx)

	// --- Cableado del módulo demo --------------------------------------------
	// De adentro hacia afuera. Cada línea recibe la de arriba ya construida, y
	// ninguna de ellas sabe quién la construyó: eso es la inyección por
	// constructor que documentaste en demo_application.go.
	demoRepo := demomongo.NewDemoRepository(db)         // adaptador de salida
	demoApp := application.NewDemoApplication(demoRepo) // casos de uso
	demoHandler := httpDir.NewDemoHandler(demoApp)      // adaptador de entrada

	demoRepository := demomongo.NewDemoRepository(db)
	_ = demoRepository

	//cargo el port desde internal/config.go
	port := cfg.Env.Port
	if port == "" {
		port = defaultPort
	}
	addr := ":" + port

	fmt.Println("Server running on", addr)

	//Rutas
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("Get /demos", demoHandler.GetAll)

	log.Fatal(http.ListenAndServe(addr, mux))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("everything is ok!"))
}
