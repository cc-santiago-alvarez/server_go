package domain

import "errors"

// Los errores de negocio del módulo se declaran aquí como variables, no como
// fmt.Errorf sueltos en cada sitio donde ocurren.
// # Cómo se usan en cada capa
//
// El dominio los devuelve tal cual:
//
//	return nil, ErrNameRequired
//
// La capa de aplicación los envuelve con %w para añadir contexto sin perder el
// error original, que es lo que hace que errors.Is siga funcionando más arriba:
//
//	return fmt.Errorf("demo: create: %w", err)
//
// Y la capa HTTP es la que decide el código de estado, porque el dominio no sabe
// qué es un 404:
//
//	switch {
//	case errors.Is(err, ErrNameRequired):
//		// 400 Bad Request
//	case errors.Is(err, ErrDemoNotFound):
//		// 404 Not Found
//	case errors.Is(err, ErrInvalidStatusTransition):
//		// 409 Conflict
//	}
//
// # Convención de los mensajes
//
// El prefijo "demo:" identifica el paquete de origen, de modo que al envolverse
// el error se lee la ruta completa de dónde falló:
//
//	demo: create: demo: NAME_IS_REQUIRED
//
// El texto en mayúsculas es el código estable que la capa HTTP devuelve al
// cliente.
var (
	ErrNameRequired            = errors.New("demo: NAME_IS_REQUIRED")
	ErrInvalidPrice            = errors.New("demo: INCORRECT_PRICE")
	ErrDemoNotFound            = errors.New("demo: NOT_FOUND")
	ErrInvalidStatusTransition = errors.New("demo: INVALID_STATUS_TRANSITION")
	ErrDemoAlreadyRemoved      = errors.New("demo: ALREADY_REMOVED")
	ErrDemoRemoved             = errors.New("demo: IS_REMOVED")
)
