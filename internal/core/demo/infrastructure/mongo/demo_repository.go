package mongo

import (
	// El driver oficial tambien se llama "mongo", igual que este paquete.
	// Se le pone alias para que quede claro que es la libreria y no codigo propio.
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	mongodb "server_go/internal/database/mongo"
)

const collectionName = "products_demo"

// DemoMongoRepository es el adaptador de salida: traduce entre el dominio y
// la coleccion de Mongo. Es el unico lugar del modulo que conoce el driver.
type DemoMongoRepository struct {
	collection *mongodriver.Collection
}

// NewDemoRepository recibe el cliente ya conectado y se queda con su coleccion.
func NewDemoRepository(db *mongodb.Client) *DemoMongoRepository {
	return &DemoMongoRepository{
		collection: db.Collection(collectionName),
	}
}
