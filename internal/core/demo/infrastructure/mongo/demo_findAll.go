package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"server_go/internal/core/demo/domain"
)

// FindAll devuelve todos los productos de la coleccion.
func (r *DemoRepository) FindAll(ctx context.Context) ([]domain.Demo, error) {
	// Filtro vacio = traer todo. Aqui van las condiciones cuando las necesites.
	filter := bson.M{}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo: find all products: %w", err)
	}
	defer cursor.Close(ctx)

	// make con len 0 y no 'var products []domain.Product': asi una coleccion
	// vacia se serializa como [] y no como null.
	products := make([]domain.Demo, 0)
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("mongo: decode products: %w", err)
	}

	return products, nil
}
