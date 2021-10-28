package main

import (
	"fmt"
	"os"

	"github.com/spandx/cache/cache"
	"github.com/spandx/cache/nuget"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "nuget":
			c := nuget.NewCatalog()
			cache := cache.NewCache(".index")
			c.Each(func(item nuget.Dependency) {
				cache.Write(item.Name, item.Version, []string{item.LicenseExpression})
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
