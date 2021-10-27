package cache

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"strings"
)

type DataFile struct {
	Path string
}

func (d DataFile) Insert(name string, version string, licenses []string) {
	file, _ := os.OpenFile(d.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	csv := fmt.Sprintf(`"%s","%s","%s"`, name, version, strings.Join(licenses, "-|-"))
	file.WriteString(csv)
	file.WriteString("\n")
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

func (c Cache) Write(name string, version string, licenses []string) error {
	if name == "" {
		return errors.New("Name is empty")
	}
	if version == "" {
		return errors.New("Version is empty")
	}

	dataFile := c.dataFileFor(name)
	dataFile.Insert(name, version, licenses)
	return nil
}
