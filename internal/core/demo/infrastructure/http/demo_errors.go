package http

import (
	"errors"
	"log"
	"net/http"
	"server_go/internal/core/demo/domain"
	"strings"
)

type errorResponse struct {
	Error string `json:"error"`
}

// La tabla es la traducción "error de negocio → código HTTP" que documentaste
// en domain/demo_errors.go. Vive aquí y no en el dominio porque el dominio no
// sabe qué es un 404.
//
// Es un slice y no un map porque el orden importa: se recorre de arriba abajo
// y gana la primera coincidencia.
var statusByError = []struct {
	err    error
	status int
}{
	{domain.ErrNameRequired, http.StatusBadRequest},
	{domain.ErrInvalidPrice, http.StatusBadRequest},
	{domain.ErrDemoNotFound, http.StatusNotFound},
	{domain.ErrInvalidStatusTransition, http.StatusConflict},
	{domain.ErrDemoAlreadyRemoved, http.StatusConflict},
	{domain.ErrDemoRemoved, http.StatusConflict},
}

// writeError es el único sitio del módulo que decide códigos de estado.
func writeError(w http.ResponseWriter, err error) {
	for _, m := range statusByError {
		if errors.Is(err, m.err) {
			writeJSON(w, m.status, errorResponse{Error: errorCode(m.err)})
			return
		}
	}

	log.Printf("demo: error inesperado: %v", err)
	writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "INTERNAL_ERROR"})
}

// errorCode extrae el código estable del centinela: "demo: NOT_FOUND" → "NOT_FOUND".
//
// Se aplica sobre el centinela y no sobre err, porque err viene envuelto y su
// texto sería "demo: get all: demo: NOT_FOUND".
func errorCode(err error) string {
	return strings.TrimPrefix(err.Error(), "demo: ")
}
