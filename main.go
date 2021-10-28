package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spandx/cache/cache"
	"github.com/spandx/cache/nuget"
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
			response, _ := http.Get("https://replicate.npmjs.com/registry/_all_docs")
			defer response.Body.Close()

			cache := cache.NewCache(".index", "npm")
			scanner := bufio.NewScanner(response.Body)
			for scanner.Scan() {
				line := strings.TrimSuffix(scanner.Text(), ",")
				var item Item
				json.Unmarshal([]byte(line), &item)

				response, _ := http.Get(fmt.Sprintf("https://replicate.npmjs.com/%s", item.Key))
				defer response.Body.Close()

				var dependency Dependency
				json.NewDecoder(response.Body).Decode(&dependency)

				for _, v := range dependency.Versions {
					cache.Write(v.Name, v.Version, []string{v.License})
				}
			}
			cache.Tidy()
			break
		default:
			fmt.Println("Unknown eco-system")
		}
	} else {
		fmt.Println("Please specify an eco system")
	}
}
