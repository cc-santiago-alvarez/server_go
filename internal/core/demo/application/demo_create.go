package application

import (
	"context"
	"fmt"

	"server_go/internal/core/demo/domain"
	"server_go/internal/core/demo/ports"
)

// Create es el caso de uso "crear demo".
//
// Son tres pasos y ninguno es una decisión: construir, persistir, devolver. Esta
// capa orquesta, no valida ni transforma. Si algún día aquí aparece un if sobre
// in.Name o in.Price, es que una regla se escapó del dominio.
func (s *DemoApplication) Create(ctx context.Context, in ports.CreateDemoInput) (*domain.Demo, error) {
	// 1. El dominio construye y valida a la vez. NewDemo es la única puerta de
	//    entrada a un Demo nuevo: de aquí no sale una entidad en estado inválido.
	//    Él pone el ID, el status inicial y los timestamps; el input no los trae.
	demo, err := domain.NewDemo(in.Name, in.Description, in.Price)
	if err != nil {
		return nil, fmt.Errorf("demo: create: %w", err)
	}

	// El repositorio persiste algo que ya es válido. Se pasa *demo porque el
	//    puerto recibe un valor: el adaptador no tiene por qué poder mutar la
	//    entidad, solo escribirla.
	if err := s.repo.InsertOne(ctx, *demo); err != nil {
		return nil, fmt.Errorf("demo: create: %w", err)
	}

	// 3. Se devuelve la entidad, no el input: trae el ID y el CreatedAt que
	//    acaba de generar el dominio y que quien llamó todavía no conocía.
	return demo, nil
}
