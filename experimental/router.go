package experimental

import (
	"context"
	"fmt"
	"strings"
)

type RouterHandlerFunc func(ctx context.Context, args []string)
type RouterContextWriterFunc func(ctx context.Context, value interface{}) context.Context

type Router struct {
	parent       *Router
	route        string
	next         map[string]*Router
	handlerFunc  RouterHandlerFunc
	handlerFuncs map[string]RouterHandlerFunc
	ctxPrefix    []string
	ctxWriters   []RouterContextWriterFunc
}

func NewRouter() *Router {
	return &Router{
		next:         make(map[string]*Router),
		handlerFuncs: make(map[string]RouterHandlerFunc),
	}
}

func (r *Router) AddNext(route string, router *Router) *Router {
	router.route = route
	router.parent = r
	r.next[route] = router
	return r
}

func (r *Router) SetStringContextPrefix(prefix string) *Router {
	prefixItems := strings.Split(prefix, " ")
	r.ctxPrefix = prefixItems
	for _, prefixItem := range prefixItems {
		r.ctxWriters = append(r.ctxWriters, NewRouterCtxWriter(prefixItem))
	}
	return r
}

func (r *Router) SetHandlerFunc(h RouterHandlerFunc) *Router {
	r.handlerFunc = h
	return r
}

func (r *Router) SetArgHandlerFunc(arg string, h RouterHandlerFunc) *Router {
	r.handlerFuncs[arg] = h
	return r
}

func (r *Router) ContextPrefix() string {
	if len(r.ctxPrefix) == 0 {
		return ""
	}
	prefixKeys := make([]string, 0, len(r.ctxPrefix))
	for _, val := range r.ctxPrefix {
		prefixKeys = append(prefixKeys, fmt.Sprintf("<%s>", val))
	}
	return strings.Join(prefixKeys, " ")
}

func (r *Router) FullRoute() string {
	route := r.route
	if contextPrefix := r.ContextPrefix(); contextPrefix != "" {
		route = fmt.Sprintf("%s %s", route, contextPrefix)
	}
	if r.parent != nil {
		route = fmt.Sprintf("%s %s", r.parent.route, route)
	}
	return route
}

func NewRouterCtxWriter(key string) func(ctx context.Context, value interface{}) context.Context {
	return func(ctx context.Context, value interface{}) context.Context {
		return context.WithValue(ctx, key, value)
	}
}
