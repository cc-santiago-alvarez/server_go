package http

import (
	"encoding/json"
	"log"
	"net/http"
	"server_go/internal/core/demo/domain"
	"time"
)

// Este es el contrato con el cliente de la API, igual que DemoDocument es el
// contrato con Mongo. Por eso es el único archivo de esta capa con tags json.
type Demo struct {
	ID          string    `json:"_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsRemove    bool      `json:"isRemove"`
}

// newDemoResponse traduce la entidad de negocio al DTO de salida.
//
// Status es domain.DemoStatus, no string: se convierte con String(), que es
// justo para lo que existe ese método.
func newDemoResponse(d domain.Demo) Demo {
	return Demo{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Status:      d.Status.String(),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		IsRemove:    d.IsRemove,
	}
}

// newDemoListResponse hace lo mismo con una colección.
//
// make con len 0 y no 'var out []Demo': mismo motivo que en FindAll, así una
// lista vacía se serializa como [] y no como null.
func newDemoListResponse(demos []domain.Demo) []Demo {
	out := make([]Demo, 0, len(demos))
	for _, d := range demos {
		out = append(out, newDemoResponse(d))
	}
	return out
}

// writeJSON centraliza el orden correcto de escritura: primero los headers,
// después el status, y solo entonces el cuerpo. Después de WriteHeader los
// headers ya viajaron, y cualquier Set posterior se ignora en silencio.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// Si el encode falla a mitad ya no se puede cambiar el status: la respuesta
	// está a medio enviar. Lo único honesto es dejar constancia en el log.
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("http: encode response: %v", err)
	}
}
