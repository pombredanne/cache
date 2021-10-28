package cache

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type lexicographically [][]byte

func (l lexicographically) Less(i, j int) bool { return bytes.Compare(l[i], l[j]) < 0 }
func (l lexicographically) Len() int           { return len(l) }
func (l lexicographically) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

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

func (c Cache) Tidy() {
	for _, file := range c.DataFiles {
		content, _ := ioutil.ReadFile(file.Path)

		lines := bytes.Split(content, []byte{'\n'})
		sort.Sort(lexicographically(lines))

		content = bytes.Join(lines, []byte{'\n'})
		ioutil.WriteFile(file.Path, content, 0644)
	}
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
