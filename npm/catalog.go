package npm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func (c Catalog) Each(visitor func(Version)) {
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

			var dependency Dependency
			json.NewDecoder(response.Body).Decode(&dependency)

			for _, v := range dependency.Versions {
				ch <- v
			}
		}
	}()

	for item := range ch {
		visitor(item)
	}
}
