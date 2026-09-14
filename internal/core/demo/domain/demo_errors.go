package domain

import "errors"

// Aca voy a registrar todos los mensajes de error que voy a mostrar en mi api
var (
	ErrNameRequired            = errors.New("demo: NAME_IS_REQUIRED")
	ErrInvalidPrice            = errors.New("demo: INCORRECT_PRICE")
	ErrDemoNotFound            = errors.New("demo: NOT_FOUND")
	ErrInvalidStatusTransition = errors.New("demo: INVALID_STATUS_TRANSITION")
	ErrDemoAlreadyRemoved      = errors.New("demo: IS_REMOVED")
)
