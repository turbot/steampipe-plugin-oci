/*
Package oci implements a steampipe plugin for OCI.

This plugin provides data that Steampipe uses to present foreign
tables that represent Oracle Cloud Infrastructure resources.
*/
package oci

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-oci"

// Plugin creates this (oci) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"oci_core_instance":                  tableCoreInstance(ctx),
			"oci_identity_authentication_policy": tableIdentityAuthenticationPolicy(ctx),
			"oci_identity_compartment":           tableIdentityCompartment(ctx),
			"oci_identity_group":                 tableIdentityGroup(ctx),
			"oci_identity_user":                  tableIdentityUser(ctx),
			"oci_region":                         tableIdentityRegion(ctx),
		},
	}
	return p
}
