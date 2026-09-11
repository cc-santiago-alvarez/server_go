package domain

import (
	"time"

	"github.com/google/uuid"
)

type Demo struct {
	ID          string     `json:"_id" bson:"_id"`
	Name        string     `json:"name" bson:"name"`
	Description string     `json:"description" bson:"description"`
	Price       float64    `json:"price" bson:"price"`
	Status      DemoStatus `json:"status" bson:"status"`
	OwnerID     string     `json:"ownerId" bson:"ownerId"`
	CreatedAt   time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updatedAt"`
	isRemove    bool       `json:"isRemove" bson:"isRemove"`
}

func (p *Demo) New() {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.isRemove = false
}
