package ports

// Aquí viven los comandos de los casos de uso: los datos que hay que darle al
// puerto de entrada para que haga algo.
//
// # Por qué están en ports y no en domain ni en http
//
// No son reglas de negocio, así que no son de domain: son la forma de una
// petición a un caso de uso. Y no son de http porque el handler depende de la
// interfaz DemoApplicationPorts, no de application; si el comando viviera en
// application, http tendría que importar la implementación para poder llamar al
// puerto, y la indirección dejaría de servir para nada.
//
// ports es el único paquete que las dos puntas ya miran. Un contrato son sus
// operaciones más el vocabulario en el que se expresan, y esto es el
// vocabulario.
//
// # Por qué no llevan tags json
//
// Estos tipos no viajan por el cable. Quien traduce el JSON del cliente a un
// comando es infrastructure/http (demo_request.go), igual que DemoDocument
// traduce en la frontera de Mongo. Cada capa tiene su representación y una
// función que hace el paso entre ellas.
//
// # Por qué no llevan Validation()
//
// Porque validar es trabajo de domain.NewDemo, y ya lo hace. Repetirlo aquí
// pondría las reglas en dos sitios, y además devolvería errores nuevos en vez
// de los centinelas de domain/demo_errors.go: errors.Is dejaría de casar y la
// capa HTTP respondería 500 donde debería responder 400.

// CreateDemoInput son los datos mínimos que el negocio necesita para crear un
// Demo.
//
// No lleva ID, Status, CreatedAt ni IsRemove: esos los decide NewDemo. Si
// estuvieran aquí, un cliente podría mandar su propio ID o crear un demo que
// nace ya borrado.
type CreateDemoInput struct {
	Name        string
	Description string
	Price       float64
}
