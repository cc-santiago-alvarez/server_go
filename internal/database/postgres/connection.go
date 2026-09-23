package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URI     string
	AppName string

	ConnectTimeout  time.Duration
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func (c Config) withDefaults() Config {
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10 * time.Second
	}
	if c.MaxConns == 0 {
		c.MaxConns = 10
	}
	if c.MaxConnLifetime == 0 {
		c.MaxConnLifetime = time.Hour
	}
	if c.MaxConnIdleTime == 0 {
		c.MaxConnIdleTime = 30 * time.Minute
	}
	if c.AppName == "" {
		c.AppName = "server_go"
	}
	return c
}

// Client envuelve al pool para que el resto de la aplicacion no dependa
// de *pgxpool.Pool directamente.
type Client struct {
	pool *pgxpool.Pool
}

// Connect crea el pool, verifica que el servidor responda y lo deja listo.
// Se llama una sola vez, al arrancar el proceso.
func Connect(ctx context.Context, cfg Config) (*Client, error) {
	cfg = cfg.withDefaults()

	if cfg.URI == "" {
		return nil, errors.New("postgres: connection URI is empty")
	}

	// ParseConfig lee la URI y deja una config editable: asi el DSN manda en
	// credenciales y host, y el codigo manda en el pool.
	poolCfg, err := pgxpool.ParseConfig(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("postgres: invalid URI: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout
	// application_name sale en pg_stat_activity: saber que proceso abrio cada
	// conexion vale oro cuando hay que matar una consulta colgada.
	poolCfg.ConnConfig.RuntimeParams["application_name"] = cfg.AppName

	// El pool es perezoso: NewWithConfig no abre ninguna conexion todavia.
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: unable to create the pool: %w", err)
	}

	// El Ping es lo que confirma red, credenciales y base antes de seguir arrancando.
	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: server dont response: %w", err)
	}

	return &Client{pool: pool}, nil
}

// Pool entrega el pool que necesita un adaptador de infraestructura.
func (c *Client) Pool() *pgxpool.Pool {
	return c.pool
}

// Close cierra el pool. Se llama en el apagado del proceso.
func (c *Client) Close() {
	if c == nil || c.pool == nil {
		return
	}
	c.pool.Close()
}
