package ports

import (
	"context"

	"server_go/internal/core/demo/domain"
)

type DemoApplicationPorts interface {
	GetAll(ctx context.Context) ([]domain.Demo, error)
}

type DemoRepository interface {
	FindAll(ctx context.Context) ([]domain.Demo, error)
}
