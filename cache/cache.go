package cache

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type DataFile struct {
	Path string
}

func (d DataFile) Insert(name string, version string, licenses []string) {
	joined := strings.Join(licenses, "-|-")
	fmt.Printf("%v: %v %v %v\n", d.Path, name, version, joined)

	file, _ := os.OpenFile(d.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{name, version, joined})
}

type Cache struct {
	Path      string
	DataFiles map[string]DataFile
}

func NewCache(dir string) Cache {
	cache := Cache{
		Path:      dir,
		DataFiles: map[string]DataFile{},
	}

	for i := 0; i < 256; i++ {
		key := fmt.Sprintf("%x", i)
		cache.DataFiles[key] = DataFile{
			Path: fmt.Sprintf("%s/%s/nuget", dir, key),
		}
	}
	return cache
}

func digestFor(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (c Cache) dataFileFor(name string) DataFile {
	digest := digestFor(name)
	index := fmt.Sprintf("%v", digest[0:2])
	return c.DataFiles[index]
}

func (c Cache) Write(name string, version string, licenses []string) {
	dataFile := c.dataFileFor(name)
	dataFile.Insert(name, version, licenses)
}
