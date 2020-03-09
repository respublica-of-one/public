package console

import (
	"context"
	"errors"
	"strings"
)

type Handler func(context.Context, []string) error

type Router struct {
	contextPrefix string
	next          map[string]Router
	handler       Handler
}

func NewRouter() Router {
	return Router{
		next: make(map[string]Router),
	}
}

func (r Router) AddNext(name string, router Router) Router {
	r.next[name] = router
	return r
}

func (r Router) SetHandler(handler Handler) Router {
	r.handler = handler
	return r
}

func (r Router) SetContextPrefix(prefix string) Router {
	r.contextPrefix = prefix
	return r
}

func (r Router) Resolve(ctx context.Context, args []string) error {
	if r.contextPrefix != "" {
		prefixParts := strings.Split(r.contextPrefix, " ")
		if len(args) < len(prefixParts) {
			return errors.New("not enough args for prefix")
		}
		for _, part := range prefixParts {
			ctx = context.WithValue(ctx, part, args[0])
			args = args[1:]
		}
	}
	if len(args) > 0 {
		if next, found := r.next[args[0]]; found {
			return next.Resolve(ctx, args[1:])
		}
	}
	if r.handler == nil {
		return errors.New("handler is nil")
	}
	return r.handler(ctx, args)
}
