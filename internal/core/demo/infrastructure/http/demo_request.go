package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"server_go/internal/core/demo/ports"
)

// Este es el contrato de entrada con el cliente, la contraparte de Demo en
// demo_response.go. Igual que ese fichero, es de los únicos de esta capa con
// tags json: el formato del cable se decide aquí y en ningún otro sitio.

var errInvalidBody = errors.New("demo: INVALID_BODY")

// createDemoRequest es el cuerpo esperado en POST /demos.
//
// Va en minúscula, al contrario que el DTO de salida: nadie fuera del adaptador
// HTTP necesita nombrarla. El DTO de respuesta es público solo porque aparece
// en la firma de newDemoResponse.
//
// No tiene los mismos campos que domain.Demo por casualidad: aquí solo están
// los que el cliente puede decidir. ID, Status, CreatedAt y IsRemove se omiten a
// propósito, porque son del dominio; si estuvieran, alguien podría crear un demo
// con su propio ID o ya borrado.
type createDemoRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// toInput traduce el DTO de transporte al comando del caso de uso.
//
// No valida reglas de negocio. Eso es de domain.NewDemo, que lo hace y devuelve
// los centinelas que writeError sabe mapear.
func (req createDemoRequest) toInput() ports.CreateDemoInput {
	return ports.CreateDemoInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}
}

// decodeJSON lee el cuerpo de la petición sobre dst.
//
// Es genérico a propósito (dst any): el mismo helper sirve para el request de
// update cuando llegue. Devolver un error en vez de escribir la respuesta deja
// la decisión del código de estado en writeError, que sigue siendo el único
// sitio del módulo que decide estados.
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)

	// Sin esto, un cliente que escribe "nombre" en lugar de "name" recibe un
	// 201 y un demo con el nombre vacío... salvo que NewDemo lo pare. Mejor
	// decirle que su campo no existe que adivinar qué quiso decir.
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("http: decode body: %w: %v", errInvalidBody, err)
	}

	return nil
}
