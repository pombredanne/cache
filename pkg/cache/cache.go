package cache

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type lexicographically [][]byte

func (l lexicographically) Less(i, j int) bool { return bytes.Compare(l[i], l[j]) < 0 }
func (l lexicographically) Len() int           { return len(l) }
func (l lexicographically) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

type DataFile struct {
	Path  string
	Lines []string
}

func (d DataFile) Insert(name string, version string, licenses []string) {
	csv := fmt.Sprintf(`"%s","%s","%s"`, name, version, strings.Join(licenses, "-|-"))
	d.Lines = append(d.Lines, csv)
}

type Cache struct {
	Path      string
	DataFiles map[string]DataFile
}

func NewCache(dir string, ecosystem string) Cache {
	cache := Cache{
		Path:      dir,
		DataFiles: map[string]DataFile{},
	}

	for i := 0; i < 256; i++ {
		key := fmt.Sprintf("%x", i)
		cache.DataFiles[key] = DataFile{
			Path: fmt.Sprintf("%s/%s/%s", dir, key, ecosystem),
		}
	}
	return cache
}

func (c Cache) Flush() {
	for _, file := range c.DataFiles {
		sort.Slice(file.Lines, func(i, j int) bool {
			return file.Lines[i] < file.Lines[j]
		})

		f, err := os.Create(file.Path)
		if err == nil {
			for _, line := range file.Lines {
				fmt.Fprintln(f, line)
			}
			defer f.Close()
		}
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
