package main

import (
	"fmt"

	"github.com/spandx/cache/nuget"
)

func main() {
	c := nuget.NewCatalog()
	c.Each(func(item nuget.Dependency) {
		fmt.Printf("\"%s\",\"%s\",\"%s\"\n", item.Name, item.Version, item.LicenseExpression)
	})
}
