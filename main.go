package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Dependency struct {
	Id    string `json:"@id"`
	Type  string `json:"@type"`
	Name  string `json:"id"`
	Range string `json:"range"`
}

type DependencyGroup struct {
	Id           string       `json:"@id"`
	Type         string       `json:"@type"`
	Dependencies []Dependency `json:""`
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
	Name                     string            `json:"id"`
	IsPrerelease             bool              `json:"isPrerelease"`
	LastEdited               time.Time         `json:"lastEdited"`
	Listed                   bool              `json:"listed"`
	PackageHash              string            `json:"packageHash"`
	PackageHashAlgorithm     string            `json:"packageHashAlgorithm"`
	PackageSize              int64             `json:"packageSize"`
	Published                time.Time         `json:"published"`
	ReleaseNotes             string            `json:"releaseNotes"`
	RequireLicenseAcceptance bool              `json:"requireLicenseAcceptance"`
	VerbatimVersion          string            `json:"verbatimVersion"`
	Version                  string            `json:"version"`
	DependencyGroups         []DependencyGroup `json:"dependencyGroups"`
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

type Catalogue struct {
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

func (c Catalogue) String() string {
	return fmt.Sprintf(
		"id: %v\ntype: %v\ncommit id: %v\ntimestamp: %v\ncount: %v\ncreated at: %v\ndeleted at: %v\nedited at: %v\n",
		c.Id,
		c.Type,
		c.CommitId,
		c.CommitTimeStamp,
		c.Count,
		c.LastCreated,
		c.LastDeleted,
		c.LastEdited,
	)
}

func (x CatalogPage) String() string {
	return fmt.Sprintf(
		" id: %v\n type: %v\n commit id: %v\n timestamp: %v\n count: %v\n",
		x.Id,
		x.Type,
		x.CommitId,
		x.CommitTimeStamp,
		x.Count,
	)
}

func main() {
	url := "https://api.nuget.org/v3/catalog0/index.json"
	response, _ := http.Get(url)
	defer response.Body.Close()

	var c Catalogue
	json.NewDecoder(response.Body).Decode(&c)

	fmt.Printf("%v", c.String())
	for i, item := range c.Items {
		fmt.Printf("%d\n%v", i, item.String())
		response, _ := http.Get(item.Id)
		defer response.Body.Close()

		var cpd CatalogPageData
		json.NewDecoder(response.Body).Decode(&cpd)
		for _, item := range cpd.Items {
			fmt.Printf("\t%s %s\n", item.Name, item.Version)

			response, _ := http.Get(item.Id)
			defer response.Body.Close()

			var pdd PackageDetailsData
			json.NewDecoder(response.Body).Decode(&pdd)
			fmt.Printf("%v\n", pdd)
		}
	}
}
