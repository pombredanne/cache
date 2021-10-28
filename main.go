package main

import (
	"fmt"
	"os"

	"github.com/spandx/cache/cache"
	"github.com/spandx/cache/npm"
	"github.com/spandx/cache/nuget"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "nuget":
			c := nuget.NewCatalog()
			cache := cache.NewCache(".index", "nuget")
			c.Each(func(item nuget.Dependency) {
				cache.Write(item.Name, item.Version, []string{item.LicenseExpression})
			})
			cache.Tidy()
			break
		case "npm":
			cache := cache.NewCache(".index", "npm")
			c := npm.NewCatalog()
			c.Each(func(v npm.Version) {
				cache.Write(v.Name, v.Version, []string{v.License})
			})
			cache.Tidy()
			break
		default:
			fmt.Println("Unknown eco-system")
		}
	} else {
		fmt.Println("Please specify an eco system")
	}
}
