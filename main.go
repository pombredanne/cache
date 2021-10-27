package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

/*
1. fetch service metadata from:
https://api.nuget.org/v3/index.json

2. Parse the RegistrationBaseUrl
    {
      "@id": "https://api.nuget.org/v3/registration5-semver1/",
      "@type": "RegistrationsBaseUrl",
      "comment": "Base URL of Azure storage where NuGet package registration info is stored"
    },
3. Make a request to base_url/{LOWER_ID}/index.json

e.g. EgoPDF.MarcoRegueira"

https://api.nuget.org/v3/registration5-semver1/egopdf.marcoregueira/index.json
*/

type PackageItemDetails struct {
	Id              string             `json:"@id"`
	Type            string             `json:"@type"`
	CommitId        string             `json:"commitId"`
	CommitTimeStamp time.Time          `json:"commitTimeStamp"`
	Entry           PackageDetailsData `json:"catalogEntry"`
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

type Dependency struct {
	Id    string `json:"@id"`
	Type  string `json:"@type"`
	Name  string `json:"id"`
	Range string `json:"range"`
}

type DependencyGroup struct {
	Id              string       `json:"@id"`
	Type            string       `json:"@type"`
	Dependencies    []Dependency `json:"dependencies"`
	TargetFramework string       `json:"targetFramework"`
}

type PackageDetailsData struct {
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
type Catalog struct {
	Id              string        `json:"@id"`
	Type            []string      `json:"@type"`
	CommitId        string        `json:"commitId"`
	CommitTimeStamp time.Time     `json:"commitTimeStamp"`
	Count           int           `json:"count"`
	LastCreated     time.Time     `json:"nuget:lastCreated"`
	LastDeleted     time.Time     `json:"nuget:lastDeleted"`
	LastEdited      time.Time     `json:"nuget:lastEdited"`
	Items           []CatalogPage `json:"items"`
}

func NewCatalog() Catalog {
	response, _ := http.Get("https://api.nuget.org/v3/catalog0/index.json")
	defer response.Body.Close()

	var c Catalog
	json.NewDecoder(response.Body).Decode(&c)
	return c
}

type Visitor func(PackageDetailsData)

func (c CatalogPage) Each(v func(PackageDetails)) {
	response, _ := http.Get(c.Id)
	defer response.Body.Close()

	var cpd CatalogPageData
	json.NewDecoder(response.Body).Decode(&cpd)

	for _, item := range cpd.Items {
		v(item)
	}
}

func (c PackageDetails) Each(v func(PackageItemDetails)) {
	registrationBaseUrl := "https://api.nuget.org/v3/registration5-semver1/"
	response, _ := http.Get(
		fmt.Sprintf("%s%s/index.json", registrationBaseUrl, strings.ToLower(c.Name)),
	)
	defer response.Body.Close()

	var data PackageIndexData
	json.NewDecoder(response.Body).Decode(&data)

	for _, item := range data.Items {
		for _, itemx := range item.Items {
			v(itemx)
		}
	}
}

func (c Catalog) Each(visitor Visitor) {
	for _, item := range c.Items {
		item.Each(func(item PackageDetails) {
			item.Each(func(x PackageItemDetails) {
				visitor(x.Entry)
			})
		})
	}
}

func main() {
	c := NewCatalog()
	ch := make(chan PackageDetailsData)
	go func() {
		c.Each(func(itemx PackageDetailsData) {
			ch <- itemx
		})
		close(ch)
	}()
	for item := range ch {
		fmt.Printf("\"%s\",\"%s\",\"%s\"\n", item.Name, item.Version, item.LicenseExpression)
	}
}
