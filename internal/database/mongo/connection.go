package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	URI      string
	Database string
	AppName  string

	ConnectTimeout time.Duration
	Timeout        time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
}

func (c Config) withDefaults() Config {
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10 * time.Second
	}
	if c.Timeout == 0 {
		c.Timeout = 15 * time.Second
	}
	if c.MaxPoolSize == 0 {
		c.MaxPoolSize = 100
	}
	if c.AppName == "" {
		c.AppName = "server_go"
	}
	return c
}

// Client envuelve al cliente del driver para que el resto de la aplicacion
// no dependa de *mongo.Client directamente.
type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

// Connect crea el pool, verifica que el servidor responda y deja lista la
// base de datos de trabajo. Se llama una sola vez, al arrancar el proceso.
func Connect(ctx context.Context, cfg Config) (*Client, error) {
	cfg = cfg.withDefaults()

	if cfg.URI == "" {
		return nil, errors.New("mongo: Connection URI is emty")
	}

	if cfg.Database == "" {
		return nil, errors.New("mongo: database name is missing")
	}

	opts := options.Client().
		ApplyURI(cfg.URI).
		SetAppName(cfg.AppName).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetServerSelectionTimeout(cfg.ConnectTimeout).
		SetTimeout(cfg.Timeout).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("mongo: Unable to create the client:%w", err)
	}

	// El Ping es lo que confirma red, credenciales y authSource antes de seguir arrancando.
	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo: server dont response: %w", err)
	}

	return &Client{
		client: client,
		db:     client.Database(cfg.Database),
	}, nil
}

// Collection entrega la coleccion que necesita un adaptador de infraestructura.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

// DB expone la base de datos para casos puntuales (indices, transacciones).
func (c *Client) DB() *mongo.Database {
	return c.db
}

// Close cierra el pool. Se llama en el apagado del proceso.
func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Disconnect(ctx)
}
