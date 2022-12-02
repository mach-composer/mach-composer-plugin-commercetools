package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-commercetools/internal"
)

func main() {
	p := internal.NewCommercetoolsPlugin()
	plugin.ServePlugin(p)
}
