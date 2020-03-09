package main

import (
	"context"
	"fmt"
	"github.com/respublica-of-one/public/pkg/console"
)

func main() {

	router := console.NewRouter().
		AddNext("first", console.NewRouter()).
		AddNext("second",
			console.NewRouter().
				SetHandler(func(ctx context.Context, strings []string) error {
					fmt.Printf("Context: %+v\nArgs: %+v\n", ctx, strings)
					return nil
				})).
		SetHandler(func(ctx context.Context, strings []string) error {
			fmt.Printf("Context: %+v\nArgs: %+v\n", ctx, strings)
			return nil
		})

	router = router.AddNext("ctx", console.NewRouter().
		SetContextPrefix("one two").
		SetHandler(func(ctx context.Context, strings []string) error {
			val1 := ctx.Value("one").(string)
			fmt.Printf("Context: %+v\nArgs: %+v\n", ctx, strings)
			fmt.Printf("one: %s\n", val1)
			return nil
		}))

	fmt.Println(router.Resolve(context.Background(), []string{"ctx", "no1", "no2", "arg1", "arg2"}))

}
