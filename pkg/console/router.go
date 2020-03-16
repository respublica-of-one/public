package console

import (
	"os"
	"strings"
)

type Router struct {
	path    string
	match   RouterMatchHandlerFunc
	next    []*Router
	resolve RouterResolver
}

func NewRouter(path string) *Router {
	return &Router{path: path}
}
func NewAppRouter() *Router {
	return &Router{
		path:    strings.ToLower(os.Args[0]),
		match:   RouterMatchIgnoreCaseStringHandler(os.Args[0]),
		next:    nil,
		resolve: nil,
	}
}
func (r *Router) SetMatcherFunc(fn RouterMatchHandlerFunc) *Router {
	r.match = fn
	return r
}
func (r *Router) AddNextRouter(router *Router) *Router {
	r.next = append(r.next, router)
	return r
}

// TODO think about the naming
func (r *Router) CreateNextRouter(router *Router) *Router {
	r.next = append(r.next, router)
	return router
}
func (r *Router) SetResolver(resolve RouterResolver) *Router {
	r.resolve = resolve
	return r
}
