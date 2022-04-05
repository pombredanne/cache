package main

import (
	"fmt"
	"os"

	"github.com/spandx/cache/pkg/cache"
	"github.com/spandx/cache/pkg/core"
	"github.com/spandx/cache/pkg/npm"
	"github.com/spandx/cache/pkg/nuget"
)

func main() {
	if len(os.Args) == 2 {
		catalogs := map[string]core.Catalog{
			"nuget": nuget.NewCatalog(),
			"npm":   npm.NewCatalog(),
		}

		ecosystem := os.Args[1]
		catalog := catalogs[ecosystem]
		if catalog == nil {
			fmt.Println("Unknown ecosystem")
		} else {
			cache := cache.NewCache(".index", ecosystem)
			catalog.Each(func(item core.Dependency) {
				cache.Write(item.Name, item.Version, item.Licenses)
			})
			cache.Flush()
		}
	} else {
		fmt.Println("Usage:")
		fmt.Println("\tgo run main.go ecosystem")
	}
}
