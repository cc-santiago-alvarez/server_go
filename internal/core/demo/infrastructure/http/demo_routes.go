package http

import "net/http"

// RegisterRoutes es el mapa de rutas del módulo demo.
//
// Vive aquí y no en main por la misma razón que writeError vive en esta capa:
// el detalle "el caso de uso GetAll se expone como GET /demos" es una decisión
// del adaptador HTTP, no del ensamblador de la aplicación. main compone
// módulos; cada módulo declara su propia superficie.
//
// RegisterRoutes Recibe la interfaz y no *http.ServeMux para no atarse a la implementación:
// si mañana el mux se envuelve en middleware o se cambia por un router de
// terceros, basta con que lo que llegue sepa registrar handlers.
func (h *DemoHandler) RegisterRoutes(mux Router) {
	// El método va en MAYÚSCULAS: ServeMux compara el método tal cual, así que
	// "Get /demos" no casaría nunca con una petición real, que llega como GET.
	mux.HandleFunc("GET /demos", h.GetAll)
}

// Router es lo mínimo que el módulo necesita de quien lo monta.
type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
