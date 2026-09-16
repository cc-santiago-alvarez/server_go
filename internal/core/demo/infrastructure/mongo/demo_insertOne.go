package mongo

import (
	"context"
	"fmt"

	"server_go/internal/core/demo/domain"
)

// InsertOne persiste un Demo que ya viene construido y validado por el dominio.
//
// Por eso el cuerpo es tan corto y no hay una sola comprobación: cuando la
// entidad llega aquí, su ID, su status y sus timestamps ya existen. Este
// adaptador no decide nada, solo traduce y escribe.
//
// Devuelve solo error porque no hay nada que informar de vuelta: el ID lo generó
// el dominio, no la base de datos. Si Mongo generase el _id, esta firma tendría
// que devolverlo y el diseño sería otro.
func (r *DemoRepository) InsertOne(ctx context.Context, demo domain.Demo) error {
	// fromDomain es la frontera: a partir de esta línea se habla en el esquema
	// de Mongo y no en el del negocio. Nunca se inserta domain.Demo directamente,
	// porque entonces los nombres de los campos de la colección dependerían de
	// los nombres de los campos de la entidad.
	if _, err := r.collection.InsertOne(ctx, fromDomain(demo)); err != nil {
		return fmt.Errorf("mongo: insert demo: %w", err)
	}

	return nil
}
