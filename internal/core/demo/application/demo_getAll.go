package application

import (
	"context"
	"fmt"
	"server_go/internal/core/demo/domain"
)

// GetAll es el caso de uso "listar demos".
// Cuando no hay registros devuelve un slice vacío, nunca nil. Ese contrato está
// declarado en el puerto y lo cumple el adaptador; aquí solo se propaga.
func (s *DemoApplication) GetAll(ctx context.Context) ([]domain.Demo, error) {
	demos, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("demo: get all: %w", err)
	}

	return demos, nil
}
