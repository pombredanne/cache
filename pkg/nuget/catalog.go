package nuget

import (
	"context"
	"fmt"
	"strings"

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
}

func NewCatalog(ctx context.Context) *Catalog {
	return &Catalog{
		ctx: ctx,
	}
}

func (catalog *Catalog) Each(visitor core.Visitor[*core.Dependency]) {
	const registrationBaseUrl = "https://api.nuget.org/v3/registration5-semver1"

	for _, c := range fetch[*CatalogData](catalog.ctx, "https://api.nuget.org/v3/catalog0/index.json").Items {
		for _, item := range fetch[*CatalogPageData](catalog.ctx, c.Id).Items {
			url := fmt.Sprintf("%s/%s/index.json", registrationBaseUrl, strings.ToLower(item.Name))
			for _, x := range fetch[*PackageIndexData](catalog.ctx, url).Items {
				for _, y := range x.Items {
					visitor(&core.Dependency{
						Name:     y.Entry.Name,
						Version:  y.Entry.Version,
						Licenses: []string{y.Entry.LicenseExpression},
					})
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
	}
	x.Check(err)
	return item
}
