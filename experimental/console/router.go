package console

import (
	"context"
	"errors"
	"strings"
)

type Handler func(context.Context, []string) error

type Router struct {
	contextPrefix   string
	next            map[string]Router
	handler         Handler
	childHandlers   map[string]Handler
	fallbackHandler Handler
}

func NewRouter() Router {
	return Router{
		next:          make(map[string]Router),
		childHandlers: make(map[string]Handler),
	}
}

func (r Router) AddNext(name string, router Router) Router {
	if router.fallbackHandler == nil && r.fallbackHandler != nil {
		router.fallbackHandler = r.fallbackHandler
	}
	names := strings.Split(name, "|")
	for _, currentName := range names {
		r.next[currentName] = router
	}
	return r
}

func (r Router) SetChildHandler(name string, handler Handler) Router {
	r.childHandlers[name] = handler
	return r
}

func (r Router) SetHandler(handler Handler) Router {
	r.handler = handler
	return r
}

func (r Router) SetFallbackHandler(handler Handler) Router {
	r.fallbackHandler = handler
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
		for key, handler := range r.childHandlers {
			names := strings.Split(key, "|")
			for _, name := range names {
				if name == args[0] {
					return handler(ctx, args[1:])
				}
			}
		}
	}
	if r.handler != nil {
		return r.handler(ctx, args)
	}
	if r.handler == nil && r.fallbackHandler == nil {
		return errors.New("handler is nil")
	}
	if r.fallbackHandler == nil {
		return errors.New("fallback handler is nil")
	}
	return r.fallbackHandler(ctx, args)
}
