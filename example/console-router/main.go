package main

import (
	"context"
	"fmt"
	"github.com/respublica-of-one/public/pkg/console"
	"strings"
)

func echo(ctx context.Context, args []string) error {
	fmt.Println("echo")
	if value := ctx.Value("value"); value != nil {
		fmt.Printf("VALUE: %+v\n", value)
	}
	path := ctx.Value(console.ROUTER_META_CTX_NAME)
	if path != nil {
		fmt.Printf("RUNNING WITH META: %+v\n", path)
	}
	fmt.Printf("\targs: %+v\n", args)
	return nil
}

func main() {

	console.RegisterNextBuilderFunc(console.RouterNextMatchPrefix("ctx:"), console.RouterMatchCtxHandler)
	console.RegisterNextBuilderFunc(console.RouterNextMatchDelimitedList(","), console.RouterMatchDelimitedListItem)

	router := console.NewRouter("application").SetMatcherFunc(console.RouterMatchIgnoreCaseStringHandler("application"))
	router.CreateNext("id list,ls ctx:value").SetHandlerFunc(echo)
	router.AddNext("id get ctx:name").
		CreateNext("id set ctx:name").AddHandlerFunc("handler one", echo)

	args := strings.Split("appLication id ls something", " ")

	resolve, err := router.Resolve(&console.RouterContext{
		Ctx:  context.Background(),
		Args: args,
	})
	if err != nil {
		fmt.Printf("error on resolve: %s\n", err)
	}
	fmt.Printf("resolve: %+v\n", resolve)
	fmt.Printf("run: %s\n", resolve.Execute())
}
