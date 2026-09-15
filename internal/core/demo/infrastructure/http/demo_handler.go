package http

import "server_go/internal/core/demo/ports"

type DemoHandler struct {
	app ports.DemoApplicationPorts
}

func NewDemoHandler(app ports.DemoApplicationPorts) *DemoHandler {
	return &DemoHandler{app: app}
}
