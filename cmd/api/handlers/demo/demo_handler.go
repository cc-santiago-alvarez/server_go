package demo

import "github.com/codecraftkit/nexus"

type DemoHandler struct {
	DemoApplication
}

func NewProductsHandler(server *nexus.Server, productApplication product_ports.ProductApplicationPorts) {
	handler := ProductsHandler{
		productApplication: productApplication,
	}

	server.Group("/products", []nexus.Endpoint{
		{Path: "GET /", HandlerFunc: handler.GetAll},
	})
}
