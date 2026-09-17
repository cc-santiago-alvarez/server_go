package ports

import (
	"context"

	"server_go/internal/core/demo/domain"
)

type DemoApplicationPorts interface {
	GetAll(ctx context.Context) ([]domain.Demo, error)
	Create(ctx context.Context, in CreateDemoInput) (*domain.Demo, error)
}

type DemoRepository interface {
	FindAll(ctx context.Context) ([]domain.Demo, error)
	InsertOne(ctx context.Context, demo domain.Demo) error
}
