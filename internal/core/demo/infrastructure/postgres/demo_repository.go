package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"server_go/internal/core/demo/ports"
	postgresdb "server_go/internal/database/postgres"
)

const tableName = "products_demo"

type DemoRepository struct {
	pool *pgxpool.Pool
}

// Si falta algun metodo del puerto, el build falla aqui y no en main.
var _ ports.DemoRepository = (*DemoRepository)(nil)

func NewDemoRepository(db *postgresdb.Client) *DemoRepository {
	return &DemoRepository{pool: db.Pool()}
}
