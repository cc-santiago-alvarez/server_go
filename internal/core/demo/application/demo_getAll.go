package application

import (
	"context"
	"fmt"
	"server_go/internal/core/demo/domain"
)

// GetAll es el caso de uso "listar productos".
func (s *DemoApplication) GetAll(ctx context.Context) ([]domain.Demo, error) {
	products, err := s.DemoRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("demo: get all products: %w", err)
	}

	return products, nil
}
