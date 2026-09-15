package application

import (
	"server_go/internal/core/demo/ports"
)

// DemoApplication es el servicio de aplicacion del modulo demo: orquesta los
// casos de uso y depende del puerto de salida, nunca del driver de Mongo.
type DemoApplication struct {
	// repo es el puerto de salida. El tipo es la interfaz, no el adaptador,
	// y ese es justo el punto: aquí no se sabe si detrás hay Mongo, Postgres
	// o un stub de test.
	repo ports.DemoRepository
}

// Verificacion en tiempo de compilacion: si a DemoApplication le falta algun
// caso de uso del puerto de entrada, el build falla aqui y no en main.
var _ ports.DemoApplicationPorts = (*DemoApplication)(nil)

// NewDemoApplication recibe las dependencias ya construidas. Es inyección de
// dependencias por constructor: el servicio no busca lo que necesita, se lo dan.
func NewDemoApplication(repo ports.DemoRepository) *DemoApplication {
	return &DemoApplication{
		repo: repo,
	}
}
