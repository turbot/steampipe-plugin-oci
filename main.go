package main

import (
	"github.com/turbot/steampipe-plugin-oci/oci"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: oci.Plugin})
}
