package core

import "fmt"

type Dependency struct {
	Name     string
	Version  string
	Licenses []string
}

type Visitor[T any] func(T)

type Catalog interface {
	Each(Visitor[*Dependency])
}

func (d *Dependency) String() string {
	return fmt.Sprintf("%s %s %v", d.Name, d.Version, d.Licenses)
}
