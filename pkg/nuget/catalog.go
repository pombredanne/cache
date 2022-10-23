package nuget

import (
	"context"
	"fmt"

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
	for _, c := range fetch[*CatalogData](catalog.ctx, catalog.url).Items {
		for _, item := range fetch[*CatalogPageData](catalog.ctx, c.Id).Items {
			for _, x := range fetch[*PackageIndexData](catalog.ctx, item.URL()).Items {
				for _, y := range x.Items {
					visitor(y.Entry.ToDependency())
				}
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
