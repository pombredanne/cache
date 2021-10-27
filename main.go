package main

import (
	"github.com/spandx/cache/cache"
	"github.com/spandx/cache/nuget"
)

func main() {
	c := nuget.NewCatalog()
	cache := cache.NewCache(".")
	c.Each(func(item nuget.Dependency) {
		cache.Write(item.Name, item.Version, item.LicenseExpression)
	})
}
