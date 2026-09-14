package application

import (
	"server_go/internal/core/demo/ports"
)

// DemoApplication es el servicio de aplicacion del modulo demo: orquesta los
// casos de uso y depende del puerto de salida, nunca del driver de Mongo.
type DemoApplication struct {
	DemoMongoRepository ports.DemoRepository
}

// Verificacion en tiempo de compilacion: si a DemoApplication le falta algun
// caso de uso del puerto de entrada, el build falla aqui y no en main.
var _ ports.DemoApplicationPorts = (*DemoApplication)(nil)

// NewDemoApplication recibe la implementacion del repositorio ya construida
// (inyeccion de dependencias desde main).
func NewDemoApplication(repo ports.DemoRepository) *DemoApplication {
	return &DemoApplication{
		DemoMongoRepository: repo,
	}
}
