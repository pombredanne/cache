package main

import (
	"fmt"
	"os"

	"github.com/spandx/cache/cache"
	"github.com/spandx/cache/core"
	"github.com/spandx/cache/npm"
	"github.com/spandx/cache/nuget"
)

func main() {
	if len(os.Args) == 2 {
		items := map[string]core.Catalog{
			"nuget": nuget.NewCatalog(),
			"npm":   npm.NewCatalog(),
		}

		catalog := items[os.Args[1]]
		if catalog == nil {
			fmt.Println("Unknown eco-system")
		} else {
			cache := cache.NewCache(".index", "nuget")
			catalog.Each(func(item core.Dependency) {
				cache.Write(item.Name, item.Version, item.Licenses)
			})
			cache.Tidy()
		}
	} else {
		fmt.Println("Please specify an eco system")
	}
}
