package internal

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
)

//go:embed templates/*
var templates embed.FS

func renderResources(cfg *SiteConfig, version string) (string, error) {
	tpl, err := template.New("main.tf.gtpl").
		Funcs(helpers.TemplateFuncs()).
		Funcs(map[string]any{
			"RenderScopes": renderScope,
			"derefBool": func(i *bool, defaultValue bool) bool { if i == nil{return defaultValue } else{ return *i} },
		}).
		ParseFS(templates, "templates/*.tf.gtpl")
	if err != nil {
		return "", err
	}

	tplContext := struct {
		Config          *SiteConfig
		ProviderVersion string
	}{
		Config:          cfg,
		ProviderVersion: version,
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, tplContext); err != nil {
		return "", err
	}
	return content.String(), nil
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
