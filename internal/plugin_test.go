package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_terraformRenderComponentVars(t *testing.T) {
	tests := []struct {
		name         string
		cfg          *SiteConfig
		componentCfg *ComponentConfig
		contains     string
		wantErr      bool
	}{
		{
			name: "Escape secrets",
			cfg: &SiteConfig{
				Stores: []CommercetoolsStore{
					{
						Key: "store-x",
					},
				},
			},
			componentCfg: &ComponentConfig{
				StoreVariables: map[string]any{
					"store-x": map[string]any{
						"my-value": "${data.sops_external.variables.data[\"mms.adyen.eu_api_key\"]}",
					},
				},
			},
			contains: `my-value = data.sops_external.variables.data["mms.adyen.eu_api_key"]`,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := terraformRenderComponentVars(tt.cfg, tt.componentCfg)
			require.NoError(t, err)

			assert.Contains(t, result, tt.contains)
		})
	}
}
