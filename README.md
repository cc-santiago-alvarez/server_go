 # Orden en que se escribe el código

 Siempre de dentro hacia fuera, para que cada capa se apoye en algo ya
 definido:

  1. domain: las reglas y la entidad (esta carpeta).
  2. ports: las interfaces, es decir qué necesita el negocio del mundo exterior.
  3. application: los casos de uso, que orquestan la entidad y los puertos.
  4. infrastructure: los adaptadores concretos, mongo para persistencia y http
     para entrega, que implementan los puertos.
  5. cmd/api/main.go: el cableado, que crea las implementaciones y las inyecta.