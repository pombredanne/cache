package nuget

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type PackageItemDetails struct {
	Id              string      `json:"@id"`
	Type            string      `json:"@type"`
	CommitId        uuid.UUID   `json:"commitId"`
	CommitTimeStamp time.Time   `json:"commitTimeStamp"`
	Entry           *Dependency `json:"catalogEntry"`
	PackageContent  string      `json:"packageContent"`
	Registration    string      `json:"registration"`
}

type PackageItem struct {
	Id              string                `json:"@id"`
	Type            string                `json:"@type"`
	CommitId        uuid.UUID             `json:"commitId"`
	CommitTimeStamp time.Time             `json:"commitTimeStamp"`
	Count           uint64                `json:"count"`
	Items           []*PackageItemDetails `json:"items"`
}

type PackageIndexData struct {
	Id              string         `json:"@id"`
	Type            []string       `json:"@type"`
	CommitId        uuid.UUID      `json:"commitId"`
	CommitTimeStamp time.Time      `json:"commitTimeStamp"`
	Count           uint64         `json:"count"`
	Items           []*PackageItem `json:"items"`
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
	Id                       string    `json:"@id"`
	Type                     string    `json:"@type"`
	Authors                  string    `json:"authors"`
	Description              string    `json:"description"`
	IconUrl                  string    `json:"iconUrl"`
	Name                     string    `json:"id"`
	Language                 string    `json:"language"`
	LicenseExpression        string    `json:"licenseExpression"`
	LicenseUrl               string    `json:"licenseUrl"`
	Listed                   bool      `json:"listed"`
	MinClientVersion         string    `json:"minClientVersion"`
	PackageContent           string    `json:"packageContent"`
	ProjectUrl               string    `json:"projectUrl"`
	Published                time.Time `json:"published"`
	RequireLicenseAcceptance bool      `json:"requireLicenseAcceptance"`
	Summary                  string    `json:"summary"`
	Tags                     []string  `json:"tags"`
	Title                    string    `json:"title"`
	Version                  string    `json:"version"`

	CommitId             string             `json:"catalog:commitId"`
	CommitTimeStamp      time.Time          `json:"catalog:commitTimeStamp"`
	Copyright            string             `json:"copyright"`
	Created              time.Time          `json:"created"`
	IsPrerelease         bool               `json:"isPrerelease"`
	LastEdited           time.Time          `json:"lastEdited"`
	PackageHash          string             `json:"packageHash"`
	PackageHashAlgorithm string             `json:"packageHashAlgorithm"`
	PackageSize          uint64             `json:"packageSize"`
	ReleaseNotes         string             `json:"releaseNotes"`
	VerbatimVersion      string             `json:"verbatimVersion"`
	DependencyGroups     []*DependencyGroup `json:"dependencyGroups"`
	Vulnerabilities      []*Vulnerability   `json:"vulnerabilities"`
}

type PackageDetails struct {
	Id              string    `json:"@id"`
	Type            string    `json:"@type"`
	CommitTimeStamp time.Time `json:"commitTimeStamp"`
	Name            string    `json:"nuget:id"`
	Version         string    `json:"nuget:version"`
	CommitId        uuid.UUID `json:"commitId"`
}

func (p *PackageDetails) URL() string {
	return fmt.Sprintf("https://api.nuget.org/v3/registration5-semver1/%s/index.json", strings.ToLower(p.Name))
}

type CatalogPageData struct {
	Id              string                 `json:"@id"`
	Type            string                 `json:"@type"`
	CommitId        uuid.UUID              `json:"commitId"`
	CommitTimeStamp time.Time              `json:"commitTimeStamp"`
	Count           int                    `json:"count"`
	Items           []*PackageDetails      `json:"items"`
	Parent          string                 `json:"parent"`
	Context         map[string]interface{} `json:"@context"`
}

type CatalogPage struct {
	Id              string    `json:"@id"`
	Type            string    `json:"@type"`
	CommitId        uuid.UUID `json:"commitId"`
	CommitTimeStamp time.Time `json:"commitTimeStamp"`
	Count           int       `json:"count"`
}
type AppendOnlyCatalog struct{}
type Permalink struct{}

type CatalogData struct {
	Id              string         `json:"@id"`
	Type            []string       `json:"@type"`
	CommitId        uuid.UUID      `json:"commitId"`
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
