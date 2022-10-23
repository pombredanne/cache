package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spandx/cache/pkg/cache"
	"github.com/spandx/cache/pkg/core"
	"github.com/spandx/cache/pkg/npm"
	"github.com/spandx/cache/pkg/nuget"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	catalogs := map[string]core.Catalog{
		"nuget": nuget.NewCatalog(ctx),
		"npm":   npm.NewCatalog(ctx),
	}

	sigs := make(chan os.Signal, 1)
	go func() {
		<-sigs
		cancel()
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if len(os.Args) != 2 {
		fmt.Println("Please specify an ecosystem")
		os.Exit(1)
	}

	ecosystem := os.Args[1]
	if catalog, ok := catalogs[ecosystem]; ok {
		cache := cache.NewCache(".index", ecosystem)
		catalog.Each(func(item *core.Dependency) {
			cache.Write(item.Name, item.Version, item.Licenses)
		})
		cache.Flush()
		os.Exit(0)
	}
	fmt.Println("Unknown ecosystem")
	os.Exit(1)
}
