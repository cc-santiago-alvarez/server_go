package postgres

import (
	"context"
	"fmt"
	"server_go/internal/core/demo/domain"
)

func (r *DemoRepository) InsertOne(ctx context.Context, demo domain.Demo) error {
	const query = `
		INSERT INTO products_demo
			(id, name, description, price, status, created_at, updated_at, is_remove)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	row := fromDomain(demo)

	_, err := r.pool.Exec(ctx, query,
		row.ID, row.Name, row.Description, row.Price,
		row.Status, row.CreatedAt, row.UpdatedAt, row.IsRemove,
	)
	if err != nil {
		return fmt.Errorf("postgres: insert demo: %w", err)
	}
	return nil
}
