package mongo

import (
	// El driver oficial tambien se llama "mongo", igual que este paquete.
	// Se le pone alias para que quede claro que es la libreria y no codigo propio.
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"server_go/internal/core/demo/ports"
	mongodb "server_go/internal/database/mongo"
)

const collectionName = "products_demo"

// DemoRepository es el adaptador de salida: traduce entre el dominio y
// la coleccion de Mongo. Es el unico lugar del modulo que conoce el driver.
type DemoRepository struct {
	collection *mongodriver.Collection
}

// Si falta algun metodo del puerto, el build falla aqui y no en main.
var _ ports.DemoRepository = (*DemoRepository)(nil)

// NewDemoRepository recibe el cliente ya conectado y se queda con su coleccion.
// Este es el constructor
func NewDemoRepository(db *mongodb.Client) *DemoRepository {
	return &DemoRepository{
		collection: db.Collection(collectionName),
	}
}
