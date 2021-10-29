package core

type Dependency struct {
	Name     string
	Version  string
	Licenses []string
}

type Catalog interface {
	Each(func(Dependency))
}
