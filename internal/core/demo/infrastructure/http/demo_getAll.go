package http

import "net/http"

// GetAll es el adaptador de entrada del caso de uso "listar demos".
func (h *DemoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// 1. Traducir la petición → argumentos del caso de uso.
	//    En GetAll no hay nada que leer; en GetByID aquí iría r.PathValue("id").

	// 2. Llamar al caso de uso. Se pasa r.Context() y no context.Background():
	//    así, si el cliente corta la conexión, la consulta a Mongo se cancela sola.
	demos, err := h.app.GetAll(r.Context())

	// 3. Traducir el error → código de estado HTTP. El dominio no sabe qué es un
	//    500, por eso la decisión vive aquí y no más abajo.
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Traducir el resultado → JSON. domain.Demo no sale nunca por el cable:
	//    antes pasa por el DTO, que es el contrato con el cliente.
	writeJSON(w, http.StatusOK, newDemoListResponse(demos))
}
