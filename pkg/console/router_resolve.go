package console

import (
	"context"
	"errors"
)

var ErrNotEnoughArgs = errors.New("not enough args")
var ErrEndpointNotDefined = errors.New("endpoint handler not defined")

type RouterContext struct {
	Ctx  context.Context
	Args []string
}

type RouterResolver interface {
	Resolve(c *RouterContext) (RouterEndpoint, error)
}

func (r *Router) Resolve(c *RouterContext) (RouterEndpoint, error) {
	if len(c.Args) == 0 {
		return nil, ErrNotEnoughArgs
	}
	match, err := r.match(c, c.Args[0])
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, nil
	}
	c.Args = c.Args[1:]
	for _, next := range r.next {
		if endpoint, err := next.Resolve(c); endpoint != nil || err != nil {
			return endpoint, err
		}
	}
	if r.resolve != nil {
		return r.resolve.Resolve(c)
	}
	return nil, ErrEndpointNotDefined
}
