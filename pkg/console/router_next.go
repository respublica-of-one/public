package console

type RouterNextAddCreator interface {
	RouterMatchResolver
	CreateNext(route string) RouterNextAddCreator
	AddNext(route string, router RouterNextAddCreator) RouterNextAddCreator
}
