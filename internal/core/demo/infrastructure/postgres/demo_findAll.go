package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"server_go/internal/core/demo/domain"
)

// FindAll devuelve todos los productos no eliminados.
func (r *DemoRepository) FindAll(ctx context.Context) ([]domain.Demo, error) {
	const query = `
		SELECT id, name, description, price, status, created_at, updated_at, is_remove
		FROM products_demo
		WHERE is_remove = FALSE
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres: find all demos: %w", err)
	}

	// CollectRows cierra rows por su cuenta, tambien si falla el escaneo:
	// por eso aqui no hace falta el defer rows.Close() del adaptador de Mongo.
	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[DemoRow])
	if err != nil {
		return nil, fmt.Errorf("postgres: scan demos: %w", err)
	}

	// make con len 0 y no 'var demos []domain.Demo': asi una tabla vacia
	// se serializa como [] y no como null.
	demos := make([]domain.Demo, 0, len(dbRows))
	for _, row := range dbRows {
		demo, err := row.toDomain()
		if err != nil {
			return nil, err
		}
		demos = append(demos, demo)
	}

	return demos, nil
}
