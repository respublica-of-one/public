package console

import (
	"context"
	"errors"
)

type RouterResolve interface {
	Resolved() bool
	Finished() bool
	Execute() error
}

type RouterMatchResolver interface {
	Match(route string) bool
	Resolve(ctx context.Context, args []string) RouterResolve
}

func RouterUnresolved() RouterResolve {
	return &unresolved{}
}

type unresolved struct{}

func (impl *unresolved) Resolved() bool {
	return false
}
func (impl *unresolved) Finished() bool {
	return true
}
func (impl *unresolved) Execute() error {
	return errors.New("args route is unresolved, no implementation given")
}
