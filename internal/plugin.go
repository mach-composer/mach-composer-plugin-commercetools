package internal

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/hashicorp/go-hclog"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type Plugin struct {
	environment string
	provider    string
	siteConfigs map[string]*SiteConfig
}

func NewCommercetoolsPlugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:    "0.30.0",
		siteConfigs: map[string]*SiteConfig{},
	}

	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier: "commercetools",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetSiteConfig:          state.SetSiteConfig,
		SetSiteComponentConfig: state.SetSiteComponentConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) IsEnabled() bool {
	return len(p.siteConfigs) > 0
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}

	cfg := SiteConfig{
		Components: map[string]ComponentConfig{},
	}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}

	if cfg.Frontend != nil {
		hclog.Default().Warn(
			fmt.Sprintf("%s: commercetools frontend block is deprecated and will be removed soon", site),
		)
	}

	if err := defaults.Set(&cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	return nil
}

func (p *Plugin) SetSiteComponentConfig(site string, component string, data map[string]any) error {
	siteConfig := p.getSiteConfig(site)
	if siteConfig == nil {
		return nil
	}

	cfg := ComponentConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	siteConfig.Components[component] = cfg

	return nil
}

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		commercetools = {
			source = "labd/commercetools"
			version = "%s"
		}
	`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	content, err := renderResources(cfg, p.provider)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (p *Plugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}
	componentCfg := cfg.getComponentSiteConfig(component)

	vars, err := terraformRenderComponentVars(cfg, componentCfg)
	if err != nil {
		return nil, err
	}

	result := &schema.ComponentSchema{
		Variables: vars,
		DependsOn: []string{"null_resource.commercetools"},
	}
	return result, nil
}

func (p *Plugin) getSiteConfig(site string) *SiteConfig {
	cfg, ok := p.siteConfigs[site]
	if !ok {
		return nil
	}
	return cfg
}

func terraformRenderComponentVars(cfg *SiteConfig, componentCfg *ComponentConfig) (string, error) {
	templateContext := struct {
		Site      *SiteConfig
		Component *ComponentConfig
	}{
		Site:      cfg,
		Component: componentCfg,
	}

	template := `
		{{ renderProperty "ct_project_key" .Site.ProjectKey }}
		{{ renderProperty "ct_api_url" .Site.APIURL }}
		{{ renderProperty "ct_auth_url" .Site.TokenURL }}
		ct_stores = {
			{{ $component := .Component -}}
			{{ range $store := .Site.Stores -}}
				{{ $store.Key }} =  {
					{{ renderProperty "key" $store.Key }}
					variables = {
						{{  range $key, $value := index $component.StoreVariables $store.Key -}}
						{{ renderProperty $key $value -}}
						{{- end -}}
					}
					secrets = {
						{{  range $key, $value := index $component.StoreSecrets $store.Key -}}
						{{ renderProperty $key $value }}
						{{- end -}}
					}
				}
			{{ end }}
		}
	`
	return helpers.RenderGoTemplate(template, templateContext)
}
