package core

import "fmt"

type Dependency struct {
	Name     string
	Version  string
	Licenses []string
}

type Catalog interface {
	Each(func(Dependency))
}

func (d *Dependency) String() string {
	return fmt.Sprintf("%s %s %v", d.Name, d.Version, d.Licenses)
}
