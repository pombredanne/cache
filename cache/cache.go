package cache

import "fmt"

type Cache struct {
	Path string
}

func NewCache(dir string) Cache {
	return Cache{
		Path: dir,
	}
}

func (c Cache) Write(name string, version string, license string) {
	fmt.Printf("\"%s\",\"%s\",\"%s\"\n", name, version, license)
}
