package sender

import (
	"context"
	"io"
	"net/http"
)

type (
	Options struct {
		ErrorMap func(error) int
	}

	JSONOptions struct {
		Options
	}

	Controller struct {
		errorMap func(error) int
		handlers map[string]func(w http.ResponseWriter, r *http.Request)
	}
)

func NewController(
	errorMap func(error) int,
) *Controller {
	return &Controller{
		errorMap: errorMap,
	}
}

func (c *Controller) Init(mux *http.ServeMux) {
	for route, handler := range c.handlers {
		mux.HandleFunc(route, handler)
	}
}

func (c *Controller) initHandlers() {
	if c.handlers == nil {
		c.handlers = make(map[string]func(w http.ResponseWriter, r *http.Request))
	}
}

func (c *Controller) Handler(route string, handler func(context.Context, Sender)) {
	c.initHandlers()
	c.handlers[route] = Handler(handler)
}

func (c *Controller) JsonHandler(route string, handler func(context.Context, Sender) (any, error)) {
	c.initHandlers()
	c.handlers[route] = JsonHandler(c.errorMap, handler)
}

func (c *Controller) StreamHandler(route string, handler func(context.Context, Sender) (io.Reader, error)) {
	c.initHandlers()
	c.handlers[route] = StreamHandler(c.errorMap, handler)
}

func (c *Controller) TemplateHandler(route string, template *Template, handler func(context.Context, Sender) (any, error)) {
	c.initHandlers()
	c.handlers[route] = TemplateHandler(c.errorMap, template, handler)
}
