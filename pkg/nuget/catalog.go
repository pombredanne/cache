package nuget

import (
	"context"
	"fmt"
	"sync"

	"github.com/spandx/cache/pkg/core"
	"github.com/xlgmokha/x/pkg/serde"
	"github.com/xlgmokha/x/pkg/x"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

type Catalog struct {
	ctx context.Context
	url string
}

func NewCatalog(ctx context.Context) *Catalog {
	return &Catalog{
		ctx: ctx,
		url: "https://api.nuget.org/v3/catalog0/index.json",
	}
}

func (catalog *Catalog) Each(visitor core.Visitor[*core.Dependency]) {
	catalogPages := make(chan *CatalogPage)
	defer func() { close(catalogPages) }()

	packageDetails := make(chan *PackageDetails)
	defer func() { close(packageDetails) }()

	packageItems := make(chan *PackageItem)
	defer func() { close(packageItems) }()

	go func() {
		catalogData := fetch[*CatalogData](catalog.ctx, catalog.url)
		for _, catalogPage := range catalogData.Items {
			catalogPages <- catalogPage
		}
		close(catalogPages)
	}()

	go func() {
		for catalogPage := range catalogPages {
			catalogPageData := fetch[*CatalogPageData](catalog.ctx, catalogPage.Id)
			for _, packageDetail := range catalogPageData.Items {
				packageDetails <- packageDetail
			}
		}

		close(packageDetails)
	}()

	go func() {
		var wg sync.WaitGroup
		for packageDetail := range packageDetails {
			wg.Add(1)

			go func() {
				defer wg.Done()
				packageIndexData := fetch[*PackageIndexData](catalog.ctx, packageDetail.URL())
				for _, packageItem := range packageIndexData.Items {
					packageItems <- packageItem
				}
			}()
		}
		wg.Wait()
		close(packageItems)
	}()

	for {
		select {
		case <-catalog.ctx.Done():
			return
		case packageItem := <-packageItems:
			for _, packageItemDetails := range packageItem.Items {
				dependency := packageItemDetails.Entry.ToDependency()
				visitor(dependency)
			}
		}
	}
}

func fetch[T any](ctx context.Context, url string) T {
	response := x.Must(otelhttp.Get(ctx, url))
	defer response.Body.Close()

	item, err := serde.From[T](response.Body, serde.JSON)
	if err != nil {
		fmt.Printf("error: %v\n", url)
		return x.Default[T]()
	}
	return item
}
