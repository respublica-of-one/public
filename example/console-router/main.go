package main

import (
	"context"
	"fmt"
	"github.com/respublica-of-one/public/pkg/args/posix"
	"github.com/respublica-of-one/public/pkg/console"
	"strings"
)

type echoArgs struct {
	Source      string `posix_args:"source,from,src,s" posix_options:"required"`
	Destination string `posix_args:"destination,to,dst,d" posix_default:"home"`
	Flag        bool   `posix_args:"flag,f" posix_options:"flag"`
}

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
	var a echoArgs
	if err := posix.ParseArgs(&a, args); err != nil {
		return err
	}
	fmt.Printf("EchoArgs: %+v\n", a)
	return nil
}

func main() {

	console.RegisterNextBuilderFunc(console.RouterNextMatchPrefix("ctx:"), console.RouterMatchCtxHandler)
	console.RegisterNextBuilderFunc(console.RouterNextMatchDelimitedList(","), console.RouterMatchDelimitedListItem)

	router := console.NewRouter("application").SetMatcherFunc(console.RouterMatchWriteExecName())
	router.CreateNext("id list,ls ctx:value").SetHandlerFunc(echo)
	router.AddNext("id get ctx:name").
		CreateNext("id set ctx:name").AddHandlerFunc("handler one", echo)

	args := strings.Split("appLication id ls something --source where -df here", " ")

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
