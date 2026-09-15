// Package domain contiene el núcleo de negocio del módulo Demo: las reglas que
// son verdad siempre, sin importar si los datos llegan por HTTP, desde un
// comando de consola o desde un test.
//
// Su única responsabilidad es garantizar que un Demo nunca exista en un estado
// inválido. Si esa garantía se cumple aquí, ninguna capa de arriba necesita
// repetir validaciones.

// # Regla de dependencias
//
// Esta capa no importa nada del proyecto: ni Mongo, ni HTTP, ni config. Solo la
// librería estándar y uuid para generar identificadores. Las demás capas
// importan a domain; domain no importa a nadie. Por eso cambiar de base de datos
// o de framework web no toca ni una línea de esta carpeta.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// contiene la entidad Demo y nada más: sus campos, su constructor, y los métodos que leen o cambian su propio estado aplicando las reglas de negocio.
// Su funcionalidad es una sola: garantizar que un Demo nunca exista en un estado inválido. Si esa garantía se cumple, cualquier caso de uso puede usar la entidad sin repetir validaciones.

// Aca definimos el modelo, este es el primer paso para construir la api
// Demo es la entidad raíz del módulo: representa un producto del catálogo.
//
// Los campos son públicos para que los adaptadores puedan construirla y leerla,
// pero el estado no se cambia asignando campos a mano. Para eso están los
// métodos de abajo, que son los únicos que aplican las reglas y mantienen
// UpdatedAt coherente:
type Demo struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Status      DemoStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsRemove    bool
}

// NewDemo es el constructor de la entidad y la única puerta de entrada para
// crear un Demo nuevo.
func NewDemo(name, description string, price float64) (*Demo, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	if price < 0 {
		return nil, ErrInvalidPrice
	}

	now := time.Now().UTC()

	return &Demo{
		ID: uuid.New().String(), Name: name, Description: description,
		Price:     price,
		Status:    DemoStatusActive,
		CreatedAt: now, UpdatedAt: now, IsRemove: false,
	}, nil
}

// Activate deja el Demo en estado activo. Devuelve error si la transición no es
// legal; las condiciones están en setStatus.
func (d *Demo) Activate() error { return d.setStatus(DemoStatusActive) }

// Deactivate deja el Demo en estado inactivo. Devuelve error si la transición no
// es legal; las condiciones están en setStatus.
func (d *Demo) Deactivate() error { return d.setStatus(DemoStatusInactive) }

// setStatus es el embudo por el que pasan todas las transiciones de estado.
// Comprueba, en este orden:
//
//  1. Que el Demo no esté borrado, porque un registro eliminado no admite
//     cambios de estado. Si lo está, devuelve ErrDemoRemoved.
//  2. Que la transición sea legal según DemoStatus.CanTransitionTo. Si no,
//     devuelve ErrInvalidStatusTransition.
//
// Solo si pasa ambas aplica el cambio y sella UpdatedAt.
func (d *Demo) setStatus(next DemoStatus) error {
	if d.IsRemove {
		return ErrDemoRemoved
	}

	if !d.Status.CanTransitionTo(next) {
		return ErrInvalidStatusTransition
	}
	d.Status = next
	d.UpdatedAt = time.Now().UTC()
	return nil
}

// Remove aplica el borrado lógico: marca el Demo como eliminado sin sacarlo de
// la base de datos, de forma que se conserva el histórico y se puede auditar.
//
// Una vez aquí, setStatus rechazará cualquier cambio de estado posterior.
func (d *Demo) Remove() error {
	if d.IsRemove {
		return ErrDemoAlreadyRemoved
	}
	d.IsRemove = true
	d.UpdatedAt = time.Now().UTC()
	return nil
}

// # Qué NO va aquí
//
//   - Tags bson o json: son detalles de Mongo y de la API. La traducción vive en
//     infrastructure/mongo (DemoDocument) e infrastructure/http.
//   - context.Context: es una preocupación de entrada/salida, y aquí no hay.
//   - DTOs de petición y respuesta: son contratos de transporte, no de negocio.
//   - Presentación (etiquetas, colores, orden de un desplegable): al negocio no
//     le importa cómo se pinta un estado.
