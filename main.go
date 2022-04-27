package main

import (
	"github.com/turbot/steampipe-plugin-oci/oci"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: oci.Plugin})
}
