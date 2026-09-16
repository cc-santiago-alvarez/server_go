package mongo

import (
	"server_go/internal/core/demo/domain"
	"time"
)

// Este es el único lugar del proyecto que conoce el esquema de la base de datos. Por eso es el único archivo con tags bson
type DemoDocument struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Price       float64   `bson:"price"`
	Status      string    `bson:"status"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	IsRemove    bool      `bson:"isRemove"`
}

// fromDomain traduce la entidad al documento que se guarda en Mongo.
//
// Status se convierte con String() porque el documento guarda un string plano:
// domain.DemoStatus es un tipo del negocio y no tiene por qué existir en la
// base de datos. El camino inverso (string → DemoStatus) tiene que validar con
// IsValid antes de entrar al dominio, por eso no es un simple cast.
func fromDomain(d domain.Demo) DemoDocument {
	return DemoDocument{
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
