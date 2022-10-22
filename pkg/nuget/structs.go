package nuget

import "time"

type PackageItemDetails struct {
	Id              string     `json:"@id"`
	Type            string     `json:"@type"`
	CommitId        string     `json:"commitId"`
	CommitTimeStamp time.Time  `json:"commitTimeStamp"`
	Entry           Dependency `json:"catalogEntry"`
}

type PackageItem struct {
	Id              string               `json:"@id"`
	Type            string               `json:"@type"`
	CommitId        string               `json:"commitId"`
	CommitTimeStamp time.Time            `json:"commitTimeStamp"`
	Count           uint64               `json:"count"`
	Items           []PackageItemDetails `json:"items"`
}

type PackageIndexData struct {
	Id              string        `json:"@id"`
	Type            []string      `json:"@type"`
	CommitId        string        `json:"commitId"`
	CommitTimeStamp time.Time     `json:"commitTimeStamp"`
	Count           uint64        `json:"count"`
	Items           []PackageItem `json:"items"`
}

type Vulnerability struct {
	Id          string `json:"@id"`
	Type        string `json:"@type"`
	AdvisoryUrl string `json:"advisoryUrl"`
	Severity    string `json:"severity"`
}

type DependencyMetadata struct {
	Id    string `json:"@id"`
	Type  string `json:"@type"`
	Name  string `json:"id"`
	Range string `json:"range"`
}

type DependencyGroup struct {
	Id              string               `json:"@id"`
	Type            string               `json:"@type"`
	Dependencies    []DependencyMetadata `json:"dependencies"`
	TargetFramework string               `json:"targetFramework"`
}

type Dependency struct {
	Id                       string            `json:"@id"`
	Type                     []string          `json:"@type"`
	Authors                  string            `json:"authors"`
	CommitId                 string            `json:"catalog:commitId"`
	CommitTimeStamp          time.Time         `json:"catalog:commitTimeStamp"`
	Copyright                string            `json:"copyright"`
	Created                  time.Time         `json:"created"`
	Description              string            `json:"description"`
	IconUrl                  string            `json:"iconUrl"`
	Name                     string            `json:"id"`
	IsPrerelease             bool              `json:"isPrerelease"`
	Language                 string            `json:"language"`
	LicenseUrl               string            `json:"licenseUrl"`
	LicenseExpression        string            `json:"licenseExpression"`
	MinClientVersion         string            `json:"minClientVersion"`
	LastEdited               time.Time         `json:"lastEdited"`
	Listed                   bool              `json:"listed"`
	PackageHash              string            `json:"packageHash"`
	PackageHashAlgorithm     string            `json:"packageHashAlgorithm"`
	ProjectUrl               string            `json:"projectUrl"`
	Summary                  string            `json:"summary"`
	Tags                     []string          `json:"tags"`
	Title                    string            `json:"title"`
	PackageSize              uint64            `json:"packageSize"`
	Published                time.Time         `json:"published"`
	ReleaseNotes             string            `json:"releaseNotes"`
	RequireLicenseAcceptance bool              `json:"requireLicenseAcceptance"`
	VerbatimVersion          string            `json:"verbatimVersion"`
	Version                  string            `json:"version"`
	DependencyGroups         []DependencyGroup `json:"dependencyGroups"`
	Vulnerabilities          []Vulnerability   `json:"vulnerabilities"`
}

type PackageDetails struct {
	Id              string    `json:"@id"`
	Type            string    `json:"@type"`
	CommitId        string    `json:"commitId"`
	CommitTimeStamp time.Time `json:"commitTimeStamp"`
	Name            string    `json:"nuget:id"`
	Version         string    `json:"nuget:version"`
}

type CatalogPageData struct {
	Id              string           `json:"@id"`
	CommitTimeStamp time.Time        `json:"commitTimeStamp"`
	Count           int              `json:"count"`
	Parent          string           `json:"parent"`
	Items           []PackageDetails `json:"items"`
}

type CatalogPage struct {
	Id              string    `json:"@id"`
	Type            string    `json:"@type"`
	CommitId        string    `json:"commitId"`
	CommitTimeStamp time.Time `json:"commitTimeStamp"`
	Count           int       `json:"count"`
}
type AppendOnlyCatalog struct{}
type Permalink struct{}

type CatalogData struct {
	Id              string         `json:"@id"`
	Type            []string       `json:"@type"`
	CommitId        string         `json:"commitId"`
	CommitTimeStamp time.Time      `json:"commitTimeStamp"`
	Count           int            `json:"count"`
	LastCreated     time.Time      `json:"nuget:lastCreated"`
	LastDeleted     time.Time      `json:"nuget:lastDeleted"`
	LastEdited      time.Time      `json:"nuget:lastEdited"`
	Items           []*CatalogPage `json:"items"`
}

type ResourceMetadata struct {
	Id      string `json:"@id"`
	Type    string `json:"@type"`
	Comment string `json:"comment"`
}

type Context struct {
	Schema  string `json:"@vocab"`
	Comment string `json:"comment"`
}

type ServiceMetadata struct {
	Version   string             `json:"version"`
	Resources []ResourceMetadata `json:"resources"`
	Context   Context            `json:"@context"`
}
