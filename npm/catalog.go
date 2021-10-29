package npm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spandx/cache/core"
)

type Item struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type Version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	License string `json:"license"`
}

type Dependency struct {
	Name     string             `json:"name"`
	Versions map[string]Version `json:"versions"`
}

type Catalog struct {
}

func NewCatalog() Catalog {
	return Catalog{}
}

func (c Catalog) Each(visitor func(core.Dependency)) {
	ch := make(chan Version)

	go func() {
		response, _ := http.Get("https://replicate.npmjs.com/registry/_all_docs")
		defer response.Body.Close()

		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			var item Item
			json.Unmarshal([]byte(strings.TrimSuffix(scanner.Text(), ",")), &item)

			response, _ := http.Get(fmt.Sprintf("https://replicate.npmjs.com/%s", item.Key))
			defer response.Body.Close()

			var d Dependency
			json.NewDecoder(response.Body).Decode(&d)

			for _, v := range d.Versions {
				ch <- v
			}
		}
		close(ch)
	}()

	for item := range ch {
		visitor(core.Dependency{
			Name:     item.Name,
			Version:  item.Version,
			Licenses: []string{item.License},
		})
	}
}
