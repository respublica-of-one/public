package main

import (
	"context"
	"fmt"
	"github.com/respublica-of-one/public/pkg/console"
)

func echo(ctx context.Context, args []string) error {
	fmt.Printf("ECHO:\n\t%+v\n\t%+v\n", ctx, args)
	return nil
}

func idSetHandler(ctx context.Context, args []string) error {
	fmt.Printf("ID IS SET: %+v\n", ctx)
	return nil
}

func idGetHandler(ctx context.Context, args []string) error {
	fmt.Printf("ID IS GOT: %+v\n", ctx)
	return nil
}

func main() {

	router := console.NewRouter().
		AddNext("echo", console.NewRouter().
			SetContextPrefix("name value").
			SetHandler(echo)).
		AddNext("error", console.NewRouter()).
		AddNext("id", console.NewRouter().
			AddNext("set", console.NewRouter().SetHandler(idSetHandler)).
			AddNext("get", console.NewRouter().SetHandler(idGetHandler)))

	fmt.Println(router.Resolve(context.Background(), []string{"id", "set", "r1", "here"}))

}
