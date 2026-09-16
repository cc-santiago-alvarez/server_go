package http

import "net/http"

// Create es el adaptador de entrada del caso de uso "crear demo".
//
// Los pasos son los mismos que en GetAll, con uno más al principio: aquí sí hay
// petición que leer. Y como en GetAll, este fichero no toma ninguna decisión de
// negocio; solo traduce en las dos direcciones.
func (h *DemoHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Traducir la petición → argumentos del caso de uso.
	//    El error de decodeJSON ya viene envuelto sobre errInvalidBody, así que
	//    writeError lo resuelve como 400 sin que haya que decidirlo aquí.s
	var req createDemoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return

	}

	// 2. Llamar al caso de uso. Se pasa r.Context() y no context.Background():
	//    si el cliente corta la conexión, la inserción en Mongo se cancela sola.
	//
	//    Fíjate en que no hay ningún if entre el decode y esta línea: las reglas
	//    de negocio las aplica NewDemo dentro del caso de uso, no el handler.
	demo, err := h.app.Create(r.Context(), req.toInput())

	// 3. Traducir el error → código de estado. Un nombre vacío llega aquí como
	//    domain.ErrNameRequired envuelto, y la tabla de demo_errors.go lo
	//    convierte en 400. El handler no sabe qué errores existen.
	if err != nil {
		writeError(w, err)
		return
	}

	// 4. Traducir el resultado → JSON. 201 y no 200 porque se creó un recurso, y
	//    con el cuerpo dentro: el cliente no conocía el ID ni CreatedAt, los
	//    generó el dominio, así que esta respuesta es su única forma de saberlos.
	//
	//    demo es *domain.Demo y newDemoResponse recibe un valor, de ahí el *demo.
	writeJSON(w, http.StatusCreated, newDemoResponse(*demo))
}
