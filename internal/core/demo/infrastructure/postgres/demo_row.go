package postgres

import (
	"fmt"
	"time"

	"server_go/internal/core/demo/domain"
)

// Este es el unico lugar del paquete que conoce el esquema de la tabla.
// Por eso es el unico archivo con tags db.
type DemoRow struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	IsRemove    bool      `db:"is_remove"`
}

// fromDomain traduce la entidad a la fila que se guarda en Postgres.
//
// Status se convierte con String() porque la columna guarda un string plano:
// domain.DemoStatus es un tipo del negocio y no tiene por que existir en la
// base de datos.
func fromDomain(d domain.Demo) DemoRow {
	return DemoRow{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Status:      d.Status.String(),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		IsRemove:    d.IsRemove,
	}
}

// toDomain es el camino inverso, y por eso valida: DemoStatus es un string por
// debajo, asi que domain.DemoStatus(r.Status) compilaria con cualquier
// contenido. Este es el filtro de la frontera para lo que viene de la base.
func (r DemoRow) toDomain() (domain.Demo, error) {
	status := domain.DemoStatus(r.Status)
	if !status.IsValid() {
		return domain.Demo{}, fmt.Errorf("postgres: demo %s: invalid status %q", r.ID, r.Status)
	}

	return domain.Demo{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Status:      status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		IsRemove:    r.IsRemove,
	}, nil
}
