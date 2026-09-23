# server_go

API en Go construida con **arquitectura hexagonal** (puertos y adaptadores),
organizada por módulos de negocio.

Este documento explica cómo está implementada la arquitectura, por qué cada
pieza está donde está, y **en qué orden** hay que escribir el código. El módulo
`demo` es la referencia: cualquier módulo nuevo se construye copiando su forma.

---

## 1. La idea en una frase

El negocio no llama a la base de datos ni al framework web: **declara lo que
necesita en forma de interfaz**, y alguien de fuera se lo trae ya construido.

Eso tiene una consecuencia práctica que se puede comprobar: cambiar de Mongo a
Postgres, o de `net/http` a otro router, no toca ni una línea de `domain` ni de
`application`.

---

## 2. La regla de dependencias

Es la única regla que no se negocia. Las dependencias apuntan **siempre hacia
adentro**:

```
     cmd/api/main.go            ← el único que conoce todas las tecnologías
            │
            ▼
  ┌──── infrastructure ────┐
  │   http        mongo    │    ← adaptadores: hablan JSON, BSON, SQL
  │    │           │       │
  └────┼───────────┼───────┘
       │           │
       ▼           ▼
        ports                    ← las interfaces: el contrato
       │           ▲
       ▼           │
   application ────┘             ← los casos de uso: orquestan
       │
       ▼
     domain                      ← las reglas: no importa a nadie
```

Leído de abajo arriba:

- **`domain`** no importa nada del proyecto. Solo `time`, `errors` y `uuid`.
- **`application`** importa `domain` y `ports`. Nunca un driver.
- **`ports`** importa `domain` (para nombrar entidades en las firmas) y nada más.
- **`infrastructure`** importa `ports` y `domain`. Es quien conoce la tecnología.
- **`main.go`** importa todo y no contiene lógica: solo construye e inyecta.

Comprobación rápida de que la regla se respeta:

```bash
# No debe devolver nada más que comentarios:
grep -rn "mongo\|postgres\|net/http" internal/core/demo/domain/ internal/core/demo/application/
```

### Por qué `application` no puede importar `infrastructure`

Porque entonces la flecha apuntaría hacia afuera y el núcleo quedaría atado al
driver. `application` depende de `ports.DemoRepository`, que es una **interfaz**.
El adaptador de Mongo también depende de `ports`, pero desde el otro lado: la
implementa. Nadie en el centro mira hacia afuera, así que un cambio en el borde
no tiene por dónde propagarse hacia adentro.

---

## 3. Mapa de carpetas

```
cmd/
  api/main.go                       Ensamblado: construye, inyecta, arranca el servidor.

internal/
  config/
    config.go                       Carga de .env.
    environments.go                 Struct con las variables ya tipadas.

  database/
    mongo/connection.go             Pool de Mongo: Connect, Collection, Close.
    postgres/connection.go          Pool de Postgres (en curso).

  core/
    demo/                           ← un módulo de negocio completo
      domain/
        demo.go                     La entidad Demo: campos, constructor, métodos de estado.
        demo_status.go              El tipo DemoStatus y sus reglas de transición.
        demo_errors.go              Los errores centinela del negocio.
      ports/
        demo_ports.go               Las interfaces: puerto de entrada y de salida.
        demo_commands.go            Los comandos de entrada (CreateDemoInput).
      application/
        demo_application.go         El servicio: struct, constructor, dependencias.
        demo_create.go              Un archivo por caso de uso.
        demo_getAll.go
      infrastructure/
        http/                       Adaptador de ENTRADA (driving).
          demo_handler.go           Struct del handler y su constructor.
          demo_routes.go            Mapa de rutas del módulo.
          demo_create.go            Un archivo por endpoint.
          demo_getAll.go
          demo_request.go           DTO de entrada + decodeJSON.
          demo_response.go          DTO de salida + writeJSON.
          demo_errors.go            Traducción error de negocio → código HTTP.
        mongo/                      Adaptador de SALIDA (driven).
          demo_repository.go        Struct del repositorio y su constructor.
          demo_document.go          DemoDocument (tags bson) + fromDomain.
          demo_insertOne.go         Un archivo por operación.
          demo_findAll.go
```

### Dos tipos de adaptador

No son lo mismo, aunque compartan carpeta:

- **De entrada (driving)**: alguien de fuera llama a la aplicación. `http` es
  uno. **Usa** el puerto `DemoApplicationPorts`.
- **De salida (driven)**: la aplicación llama a alguien de fuera. `mongo` es
  uno. **Implementa** el puerto `DemoRepository`.

Por eso las dos flechas del diagrama apuntan a `ports` desde lados opuestos.

---

## 4. El orden en que se escribe el código

Siempre **de dentro hacia fuera**, para que cada capa se apoye en algo ya
definido. Si escribes el handler primero, acabarás inventando la firma del caso
de uso y reescribiéndola tres veces.

### Paso 1 — `domain`: la entidad y sus reglas

Empiezas por lo que es verdad siempre, sin importar si los datos llegan por
HTTP, por un comando de consola o desde un test.

Escribes, en este orden:

1. **El struct de la entidad.** Campos públicos para que los adaptadores puedan
   construirla y leerla.
2. **Los tipos propios del negocio** (`DemoStatus`), con sus validadores
   (`IsValid`) y sus reglas de transición (`CanTransitionTo`).
3. **Los errores centinela** (`demo_errors.go`), con `errors.New`. Son los que
   las capas de arriba compararán con `errors.Is`.
4. **El constructor** (`NewDemo`): la única puerta de entrada a una entidad
   nueva. Valida y devuelve error, de forma que **no puede existir un Demo en
   estado inválido**.
5. **Los métodos que cambian el estado** (`Activate`, `Remove`), que aplican las
   reglas y mantienen `UpdatedAt` coherente.

Lo que **no** va aquí: tags `bson` o `json`, `context.Context`, DTOs de
transporte, ni nada de presentación.

> **La prueba:** si el constructor puede devolver una entidad que viole una
> regla, la regla está en el sitio equivocado.

### Paso 2 — `ports`: el contrato

Ahora declaras qué necesita el negocio del mundo exterior, **sin decir quién lo
va a cumplir**.

```go
// Puerto de ENTRADA: lo que la aplicación ofrece. Lo implementa application.
type DemoApplicationPorts interface {
	GetAll(ctx context.Context) ([]domain.Demo, error)
	Create(ctx context.Context, in CreateDemoInput) (*domain.Demo, error)
}

// Puerto de SALIDA: lo que la aplicación necesita. Lo implementa infrastructure.
type DemoRepository interface {
	FindAll(ctx context.Context) ([]domain.Demo, error)
	InsertOne(ctx context.Context, demo domain.Demo) error
}
```

Y en `demo_commands.go`, los **comandos**: los datos que hay que darle a un caso
de uso para que haga algo (`CreateDemoInput`).

Tres decisiones que conviene entender:

- **Los comandos viven en `ports`, no en `application` ni en `http`.** No son
  reglas de negocio, así que no son de `domain`. Y si vivieran en `application`,
  el handler tendría que importar la implementación para poder llamar al puerto,
  con lo que la indirección dejaría de servir para nada. `ports` es el único
  paquete que las dos puntas ya miran.
- **No llevan tags `json`.** No viajan por el cable: quien traduce el JSON del
  cliente a un comando es `infrastructure/http`.
- **No llevan `Validate()`.** Validar es trabajo de `domain.NewDemo`. Repetirlo
  aquí pondría las reglas en dos sitios y devolvería errores nuevos en vez de
  los centinelas, con lo que `errors.Is` dejaría de casar y la capa HTTP
  respondería 500 donde debería responder 400.
- **El comando no trae `ID`, `Status`, `CreatedAt` ni `IsRemove`.** Esos los
  decide el dominio. Si estuvieran, un cliente podría mandar su propio ID o
  crear un demo que nace ya borrado.

### Paso 3 — `application`: los casos de uso

Primero el servicio y su constructor (`demo_application.go`):

```go
type DemoApplication struct {
	repo ports.DemoRepository   // la interfaz, no el adaptador
}

// Si falta algún caso de uso del puerto de entrada, el build falla aquí.
var _ ports.DemoApplicationPorts = (*DemoApplication)(nil)

func NewDemoApplication(repo ports.DemoRepository) *DemoApplication {
	return &DemoApplication{repo: repo}
}
```

Es **inyección por constructor**: el servicio no busca lo que necesita, se lo
dan. Por eso en un test puedes pasarle un stub sin levantar una base de datos.

Después, **un archivo por caso de uso**. Cada uno orquesta y nada más:

```go
func (s *DemoApplication) Create(ctx context.Context, in ports.CreateDemoInput) (*domain.Demo, error) {
	demo, err := domain.NewDemo(in.Name, in.Description, in.Price)  // 1. construir y validar
	if err != nil {
		return nil, fmt.Errorf("demo: create: %w", err)
	}

	if err := s.repo.InsertOne(ctx, *demo); err != nil {           // 2. persistir
		return nil, fmt.Errorf("demo: create: %w", err)
	}

	return demo, nil                                                // 3. devolver
}
```

> **La prueba:** si en un caso de uso aparece un `if` sobre `in.Name` o
> `in.Price`, una regla de negocio se escapó del dominio. Esta capa orquesta,
> no valida.

El `%w` importa: envuelve el error conservando el centinela, para que
`errors.Is` siga funcionando arriba.

### Paso 4 — `infrastructure`: los adaptadores

Ahora, y solo ahora, aparece la tecnología. Los dos lados son independientes y
se pueden escribir en cualquier orden.

#### 4a. Adaptador de salida (persistencia)

1. **El struct y su constructor** (`demo_repository.go`), con la verificación de
   compilación:

   ```go
   var _ ports.DemoRepository = (*DemoRepository)(nil)
   ```

   Esa línea es lo que hace que un método que falte se reporte **aquí**, en el
   adaptador, y no en `main` a veinte archivos de distancia.

2. **El modelo de persistencia y su traducción** (`demo_document.go`): el struct
   con tags `bson` (o `db` en Postgres) y las funciones `fromDomain` / `toDomain`.

   Este es el único archivo del paquete que conoce el esquema. Nunca se inserta
   `domain.Demo` directamente: si lo hicieras, los nombres de las columnas
   pasarían a depender de los nombres de los campos de la entidad, y renombrar
   un campo del negocio sería una migración.

   `toDomain` **valida** lo que viene de la base antes de dejarlo entrar: si
   alguien escribió `'activo'` a mano en la columna `status`, el cast
   `domain.DemoStatus(s)` compilaría igual y metería una entidad inválida en el
   negocio. Por eso pasa por `IsValid`.

3. **Un archivo por operación** (`demo_insertOne.go`, `demo_findAll.go`).

#### 4b. Adaptador de entrada (HTTP)

1. **El struct y su constructor** (`demo_handler.go`). Recibe
   `ports.DemoApplicationPorts`, la interfaz, no `*application.DemoApplication`.

2. **Los DTOs y su traducción** (`demo_request.go`, `demo_response.go`). Son los
   únicos archivos de la capa con tags `json`: el formato del cable se decide
   aquí y en ningún otro sitio.

   `createDemoRequest` no tiene los mismos campos que `domain.Demo` por
   casualidad: solo están los que el cliente puede decidir.

3. **Un archivo por endpoint** (`demo_create.go`, `demo_getAll.go`). El patrón
   es siempre el mismo: decodificar → traducir a comando → llamar al caso de uso
   → traducir la entidad a DTO → escribir.

4. **La traducción de errores** (`demo_errors.go`). `writeError` es el único
   sitio del módulo que decide códigos de estado, con una tabla
   `error de negocio → status`. Vive aquí porque **el dominio no sabe qué es un
   404**.

5. **Las rutas** (`demo_routes.go`). El módulo declara su propia superficie:

   ```go
   func (h *DemoHandler) RegisterRoutes(mux Router) {
       mux.HandleFunc("GET /demos", h.GetAll)
       mux.HandleFunc("POST /demos", h.Create)
   }
   ```

   Recibe una interfaz `Router` y no `*http.ServeMux` para no atarse a la
   implementación. El verbo va en MAYÚSCULAS: `ServeMux` compara el método tal
   cual, así que `"Get /demos"` no casaría nunca.

### Paso 5 — `cmd/api/main.go`: el cableado

El último paso, y el único lugar donde se decide qué tecnología concreta se usa.
De adentro hacia afuera, cada línea recibe la anterior ya construida:

```go
demoRepo    := demomongo.NewDemoRepository(db)      // adaptador de salida
demoApp     := application.NewDemoApplication(demoRepo)  // casos de uso
demoHandler := httpDir.NewDemoHandler(demoApp)      // adaptador de entrada
```

Ninguna de esas tres sabe quién la construyó. `main` compone módulos; cada
módulo declara sus rutas. La versión de la API se decide aquí porque vale para
todos los módulos.

---

## 5. Las fronteras de traducción

Un dato cruza **cuatro** representaciones distintas entre el cliente y la base
de datos. Cada frontera tiene una función que hace el paso, y ninguna capa
conoce la representación de las otras:

```
JSON del cliente
   │  decodeJSON
   ▼
createDemoRequest      (tags json)         infrastructure/http
   │  toInput()
   ▼
ports.CreateDemoInput  (sin tags)          ports
   │  domain.NewDemo()  ← aquí se valida
   ▼
domain.Demo            (sin tags)          domain
   │  fromDomain()
   ▼
DemoDocument / DemoRow (tags bson / db)    infrastructure/mongo|postgres
   │
   ▼
base de datos
```

Y de vuelta: `toDomain()` valida al entrar, `newDemoResponse()` traduce al salir.

Parece repetitivo, y es deliberado: es lo que permite que añadir un campo
interno no cambie el contrato de la API, y que renombrar una columna no rompa a
los clientes.

---

## 6. Convenciones del proyecto

| Convención | Por qué |
|---|---|
| **Un archivo por caso de uso / operación / endpoint** | `demo_create.go` existe en tres capas. El nombre te dice *qué* hace; la carpeta, *en qué capa*. |
| **`var _ Interfaz = (*Tipo)(nil)`** en cada implementación | Un método que falte falla en el archivo culpable, no en `main`. |
| **Nombres en `lowerCamelCase`** tras el prefijo | `demo_findAll.go`, `demo_insertOne.go`. |
| **Prefijo del módulo en cada archivo** | `demo_*.go`. Al abrir el editor sabes de qué módulo es sin mirar la ruta. |
| **Errores envueltos con `%w`** y prefijo de capa | `fmt.Errorf("mongo: insert demo: %w", err)`. El prefijo dice dónde falló; `%w` conserva el centinela para `errors.Is`. |
| **Errores centinela en `domain`**, traducción a HTTP en `infrastructure/http` | El negocio no sabe qué es un 404. |
| **Constructores `NewX` que reciben dependencias** | Inyección por constructor: testeable sin infraestructura. |
| **Interfaces pequeñas definidas por quien las consume** | `Router` solo declara `HandleFunc`, que es lo único que el módulo necesita. |

---

## 7. Recetas

### Añadir un caso de uso a un módulo existente

Ejemplo: "desactivar un demo".

1. `domain` — ¿existe ya el método? `Deactivate()` sí. Si no, se escribe ahí.
2. `ports/demo_commands.go` — añade `DeactivateDemoInput` si necesita datos.
3. `ports/demo_ports.go` — añade `Deactivate(...)` a `DemoApplicationPorts`, y
   `UpdateOne(...)` a `DemoRepository` si hace falta persistir.
4. **El build se rompe en todos los adaptadores.** Es intencionado: te garantiza
   que ningún motor se queda atrás sin que nadie lo note.
5. `application/demo_deactivate.go` — el caso de uso.
6. `infrastructure/mongo/demo_updateOne.go` — la operación.
7. `infrastructure/http/demo_deactivate.go` + la ruta en `demo_routes.go`.
8. `demo_errors.go` (http) — mapea los errores nuevos si los hay.

`main.go` **no se toca**: las firmas de los constructores no cambiaron.

### Añadir un adaptador de persistencia nuevo

**No se crea una carpeta `ports` nueva.** El puerto pertenece al núcleo y es
agnóstico: Mongo y Postgres son dos formas de cumplir el mismo contrato.

1. `internal/database/<motor>/connection.go` — pool, `Connect`, `Close`.
2. `internal/core/demo/infrastructure/<motor>/` — repositorio, modelo de
   persistencia con su `fromDomain`/`toDomain`, y un archivo por operación.
3. Cada paquete necesita **su propio** `fromDomain`: los paquetes de Go no
   comparten identificadores, y además cada uno traduce a su esquema. Esa
   duplicación es correcta; unificarla haría que el esquema de un motor mandara
   sobre el otro.
4. `main.go` — cambia la línea del constructor del repositorio. Nada más.

La única razón legítima para una interfaz nueva es que un motor ofrezca algo que
otro no puede (transacciones multi-tabla, full-text). En ese caso se añade otra
interfaz pequeña **en el mismo `ports`**, expresada en lenguaje de negocio.

### Añadir un módulo nuevo

Copia el esqueleto de `demo` y recorre los pasos 1→5. Cada módulo tiene su
propio `ports`: no hay un `ports` global. Al final, en `main.go`:

```go
userHandler.RegisterRoutes(api)
```

---

## 8. Puesta en marcha

```bash
cp .env.example .env     # y rellena los valores
go mod tidy
go run ./cmd/api
```

Variables (ver `.env.example`):

| Variable | Descripción |
|---|---|
| `PORT` | Puerto HTTP. Por defecto `3252`. |
| `MONGO_URI` | Cadena de conexión de Mongo. |
| `MONGO_DATABASE` | Base de datos de trabajo. |
| `POSTGRES_URI` | DSN de Postgres. `sslmode=disable` solo en local. |

Comprobación: `GET /health` responde `200`.
Las rutas del módulo cuelgan de `/api/v1` (`GET /api/v1/demos`,
`POST /api/v1/demos`).

### Dependencias

`go get` pide **paquetes**, no módulos. Con el *module graph pruning* de Go
1.17+, `go get github.com/jackc/pgx/v5` no resuelve las dependencias de sus
subpaquetes: hay que pedir `github.com/jackc/pgx/v5/pgxpool`. Con los imports
ya escritos, `go mod tidy` lo hace solo.

---

## 9. Estado

- **Módulo `demo`**: completo sobre Mongo (`GetAll`, `Create`).
- **Adaptador de Postgres**: en curso.
