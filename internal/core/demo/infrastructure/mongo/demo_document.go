package mongo

import "time"

// Este es el único lugar del proyecto que conoce el esquema de la base de datos. Por eso es el único archivo con tags bson
type DemoDocument struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Price       float64   `bson:"price"`
	Status      string    `bson:"status"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	IsRemove    bool      `bson:"isRemove"`
}
