package main

import (
	"context"
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

	mongoCfg := mongodb.Config{
		URI:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DATABASE"),
	}

	db, err := mongodb.Connect(ctx, mongoCfg)
	if err != nil {
		log.Fatalf("mongo: conexión fallida: %v", err)
	}
	defer db.Close(ctx)

	log.Printf("mongo: conectado a la base de datos %q", mongoCfg.Database)

	// --- Cableado del módulo demo --------------------------------------------
	// De adentro hacia afuera. Cada línea recibe la de arriba ya construida, y
	// ninguna de ellas sabe quién la construyó: eso es la inyección por
	// constructor que documentaste en demo_application.go.
	demoRepo := demomongo.NewDemoRepository(db)         // adaptador de salida
	demoApp := application.NewDemoApplication(demoRepo) // casos de uso
	demoHandler := httpDir.NewDemoHandler(demoApp)      // adaptador de entrada

	//cargo el port desde internal/config.go
	port := cfg.Env.Port
	if port == "" {
		port = defaultPort
	}
	addr := ":" + port

	log.Printf("http: servidor escuchando en %s", addr)

	// main solo monta módulos: cada uno sabe qué rutas expone y bajo qué verbo.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	demoHandler.RegisterRoutes(mux)

	log.Fatal(http.ListenAndServe(addr, mux))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("everything is ok!"))
}
