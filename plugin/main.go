package plugin

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-commercetools/internal"
)

func Serve() {
	p := internal.NewCommercetoolsPlugin()
	plugin.ServePlugin(p)
}
