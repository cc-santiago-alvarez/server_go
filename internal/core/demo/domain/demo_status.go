package domain

type DemoStatus string

const (
	DemoStatusActive   DemoStatus = "active"
	DemoStatusInactive DemoStatus = "inactive"
)

// getters
// String devuelve el estado como texto plano.
func (s DemoStatus) String() string {
	return string(s)
}

// validators
// IsActive indica si el estado es el activo. Sirve para que las capas de arriba
// pregunten por el estado sin comparar contra la constante ni contra un string.
func (s DemoStatus) IsActive() bool {
	return s == DemoStatusActive
}

// IsInactive indica si el estado es el inactivo.
func (s DemoStatus) IsInactive() bool {
	return s == DemoStatusInactive
}

// IsValid indica si el valor corresponde a un estado real del negocio.
//
// Hace falta porque DemoStatus es un string por debajo, así que una conversión
// como DemoStatus(textoDelCliente) compila con cualquier contenido. Este es el
// filtro de la frontera: todo estado que venga de fuera (JSON, base de datos)
// se valida aquí antes de entrar al dominio.
func (s DemoStatus) IsValid() bool {
	switch s {
	case DemoStatusActive, DemoStatusInactive:
		return true
	}
	return false
}

// CanTransitionTo indica si se puede pasar del estado actual a next
func (s DemoStatus) CanTransitionTo(next DemoStatus) bool {
	if !next.IsValid() {
		return false
	}
	return s != next
}

// type ListDemoStatus struct {
// 	Label string `json:"name"`
// 	Value string `json:"value"`
// 	Color string `json:"color"`
// }

// func ListStatus() []ListDemoStatus {
// 	return []ListDemoStatus{
// 		ListDemoStatus{
// 			Label: DemoStatusActive.ToString(),
// 			Value: DemoStatusInactive.ToString(),
// 		},
// 		ListDemoStatus{
// 			Label: DemoStatusActive.ToString(),
// 			Value: DemoStatusInactive.ToString(),
// 		},
// 	}
// }
