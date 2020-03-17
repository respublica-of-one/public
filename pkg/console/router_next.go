package console

import (
	"errors"
	"strings"
)

var ErrBuilderNotFound = errors.New("builder not found")

type RouterNextMatchFunc func(pattern string) (string, bool)
type RouterMatchHandlerBuilderFunc func(pattern string) RouterMatchHandlerFunc

func RouterNextMatchBool(val bool) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		return pattern, val
	}
}
func RouterNextMatchString(value string) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		if value == pattern {
			return pattern, true
		}
		return pattern, false
	}
}
func RouterNextMatchPrefix(prefix string) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		if strings.HasPrefix(pattern, prefix) {
			return strings.Replace(pattern, prefix, "", 1), true
		}
		return pattern, false
	}
}
func RouterNextContainsString(value string) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		if strings.Contains(pattern, value) {
			return pattern, true
		}
		return pattern, false
	}
}
func RouterNextMatchDelimitedList(delimiter string) RouterNextMatchFunc {
	return func(pattern string) (string, bool) {
		if !strings.Contains(pattern, delimiter) {
			return pattern, false
		}
		return strings.Join(strings.Split(pattern, delimiter), "\n"), true
	}
}

type RouterNextBuilder struct {
	match RouterNextMatchFunc
	build RouterMatchHandlerBuilderFunc
}

type RouterNextBuilders struct {
	builders []RouterNextBuilder
	fallback RouterMatchHandlerBuilderFunc
}

func (n *RouterNextBuilders) AddBuilder(match RouterNextMatchFunc, build RouterMatchHandlerBuilderFunc) {
	n.builders = append(n.builders, RouterNextBuilder{
		match: match,
		build: build,
	})
}
func (n *RouterNextBuilders) SetFallback(fn RouterMatchHandlerBuilderFunc) {
	n.fallback = fn
}
func (n *RouterNextBuilders) BuildNext(pattern string) (RouterMatchHandlerFunc, error) {
	for _, next := range n.builders {
		if key, match := next.match(pattern); match {
			return next.build(key), nil
		}
	}
	if n.fallback != nil {
		return n.fallback(pattern), nil
	}
	return nil, ErrBuilderNotFound
}

var nextBuildersInstance *RouterNextBuilders

func SetNextBuilders(n *RouterNextBuilders) {
	nextBuildersInstance = n
}
func nextBuilders() *RouterNextBuilders {
	if nextBuildersInstance == nil {
		initDefaultNextBuilders()
	}
	return nextBuildersInstance
}
func RegisterNextBuilderFunc(match RouterNextMatchFunc, fn RouterMatchHandlerBuilderFunc) {
	nextBuilders().AddBuilder(match, fn)
}
func buildNext(pattern string) (RouterMatchHandlerFunc, error) {
	return nextBuilders().BuildNext(pattern)
}
func initDefaultNextBuilders() {
	nextBuildersInstance = &RouterNextBuilders{
		builders: []RouterNextBuilder{},
		fallback: RouterMatchIgnoreCaseStringHandler,
	}
}
