package nuget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/spandx/cache/core"
	"golang.org/x/sync/semaphore"
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

func NewCatalog() Catalog {
	response, err := http.Get("https://api.nuget.org/v3/catalog0/index.json")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return Catalog{}
	}
	defer response.Body.Close()

	var c Catalog
	json.NewDecoder(response.Body).Decode(&c)
	return c
}

func (c Catalog) Each(visitor func(core.Dependency)) {
	const registrationBaseUrl = "https://api.nuget.org/v3/registration5-semver1/"
	dependencies := make(chan core.Dependency, 2048)
	lock := semaphore.NewWeighted(int64(256))
	urls := make(chan string, 1024)

	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	wg2.Add(len(c.Items))

	go func() {
		for {
			fmt.Printf("urls: %d, dependencies: %d\n", len(urls), len(dependencies))
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for _, c := range c.Items {
			go func(c CatalogPage) {
				defer wg2.Done()

				lock.Acquire(context.TODO(), 1)
				response, err := http.Get(c.Id)
				lock.Release(1)

				if err == nil {
					defer response.Body.Close()

					var cpd CatalogPageData
					json.NewDecoder(response.Body).Decode(&cpd)

					wg.Add(len(cpd.Items))
					for _, item := range cpd.Items {
						urls <- fmt.Sprintf("%s%s/index.json", registrationBaseUrl, strings.ToLower(item.Name))
					}
				} else {
					fmt.Printf("%s\n", err.Error())
				}
			}(c)
		}
	}()

	go func() {
		for url := range urls {
			go func(url string) {
				lock.Acquire(context.TODO(), 1)
				response, err := http.Get(url)
				lock.Release(1)

				if err == nil {
					defer response.Body.Close()

					var data PackageIndexData
					json.NewDecoder(response.Body).Decode(&data)

					for _, x := range data.Items {
						for _, y := range x.Items {
							dependencies <- core.Dependency{
								Name:     y.Entry.Name,
								Version:  y.Entry.Version,
								Licenses: []string{y.Entry.LicenseExpression},
							}
						}
					}
				} else {
					fmt.Printf("%s\n", err.Error())
				}

				wg.Done()
			}(url)
		}
	}()

	go func() {
		for item := range dependencies {
			visitor(item)
		}
	}()

	wg2.Wait()
	close(urls)
	wg.Wait()
	close(dependencies)
}
