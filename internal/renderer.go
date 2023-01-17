package internal

import (
	"embed"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/flosch/pongo2/v5"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
)

//go:embed templates/*
var templates embed.FS

func renderResources(cfg *SiteConfig, version string) (string, error) {
	templateSet := pongo2.NewSet("", &helpers.EmbedLoader{Content: templates})
	template := pongo2.Must(templateSet.FromFile("main.tf"))

	return template.Execute(pongo2.Context{
		"commercetools":    cfg,
		"render_scopes":    renderScope,
		"provider_version": version,
	})
}

var STORE_SUPPORTED_SCOPES = []string{
	"manage_orders",
	"manage_my_orders",
	"view_orders",
	"manage_customers",
	"view_customers",
	"manage_my_profile",
}

func renderScope(scopes []string, projectKey string, storeKey string) string {

	sl := make([]string, 0)
	for _, scope := range scopes {
		sl = append(sl, fmt.Sprintf(`"%s:%s",`, scope, projectKey))

		if storeKey != "" && pie.Contains(STORE_SUPPORTED_SCOPES, scope) {
			sl = append(sl, fmt.Sprintf(`"%s:%s:%s",`, scope, projectKey, storeKey))
		}
	}

	result := fmt.Sprintf("[\n  %s\n]", strings.Join(sl, "\n"))
	return result

}
