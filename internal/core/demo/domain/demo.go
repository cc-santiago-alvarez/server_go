package domain

import (
	"time"

	"github.com/google/uuid"
)

// contiene la entidad Demo y nada más: sus campos, su constructor, y los métodos que leen o cambian su propio estado aplicando las reglas de negocio.
// Su funcionalidad es una sola: garantizar que un Demo nunca exista en un estado inválido. Si esa garantía se cumple, cualquier caso de uso puede usar la entidad sin repetir validaciones.

// Aca definimos el modelo, este es el primer paso para construir la api
type Demo struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Status      DemoStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsRemoved   bool
}

// Este es su constructor
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
		CreatedAt: now, UpdatedAt: now, IsRemoved: false,
	}, nil
}

// Activate y Deactivate son la unica forma de cambiar el estado desde fuera.
func (d *Demo) Activate() error   { return d.setStatus(DemoStatusActive) }
func (d *Demo) Deactivate() error { return d.setStatus(DemoStatusInactive) }

// setStatus es privado: concentra la regla para que Activate y Deactivate
// no la repitan.
func (d *Demo) setStatus(next DemoStatus) error {
	if !d.Status.CanTransitionTo(next) {
		return ErrInvalidStatusTransition
	}
	d.Status = next
	d.UpdatedAt = time.Now().UTC()
	return nil
}

func (d *Demo) Remove() error {
	if d.IsRemoved {
		return ErrDemoAlreadyRemoved
	}
	d.IsRemoved = true
	d.UpdatedAt = time.Now().UTC()
	return nil
}
