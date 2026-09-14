package mongo

import "time"

type DemoDocument struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Price     float64   `bson:"price"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"createdAt"`
	IsRemoved bool      `bson:"isRemove"`
}
