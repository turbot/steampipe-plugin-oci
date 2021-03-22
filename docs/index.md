---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/oci.svg"
brand_color: "#FF9900"
display_name: "Oracle Cloud Infrastructure"
short_name: "oci"
description: "Steampipe plugin for Oracle Cloud Infrastructure services and resource types."
---

# OCI

The Oracle Cloud Infrastructure (OCI) plugin is used to interact with the many resources supported by OCI.

### Installation

To download and install the latest oci plugin:

```bash
steampipe plugin install oci
```

Installing the latest oci plugin will create a connection config file (`~/.steampipe/config/oci.spc`) with a single default connection named `oci`.

Note that there is nothing special about the default connection, other than that it is created by default on plugin install - You can delete or rename this connection, or modify its configuration options (via the configuration file).

## Connection Configuration

Connection configurations are defined using HCL in one or more Steampipe config files. Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.

### Scope

Each OCI connection is scoped to a single OCI Tenant, with a single set of credentials. You may configure multiple OCI connections if desired, with each connecting to a different tenant. Each OCI connection may be configured for multiple regions.

### Configuration Arguments

The OCI plugin allows you set static credentials with the `config_file_profile` and `config_path` arguments. You may select one or more regions with the `regions` argument.
An OCI connection may connect to multiple regions, however be aware that performance may be negatively affected by both the number of regions and the latency to them.

```hcl
# credentials via user API Key pair stored in config file
connection "oci_tenant_x" {
  plugin                = "oci"
  config_file_profile   = "DEFAULT"
  config_path           = "~/.oci/config"
  regions               = ["ap-mumbai-1" , "us-ashburn-1"]
}
```

Alternatively, you may set details required for API Key authentication using `user_ocid`, `fingerprint`, `private_key`, `private_key_path`, and `private_key_password` argument.

````hcl
# credentials via profile
connection "oci_tenant_y" {
  plugin            = "oci"
  tenancy_ocid      = "test"
  user_ocid         = "test_user"
  fingerprint       = "dummy-fingerprint"
  private_key_path  = "~/.ssh/oci_private.pem"
  regions               = ["ap-mumbai-1" , "us-ashburn-1"]
}


If no credentials are specified, the plugin will use the OCI credentials resolver to get the current credentials in the same manner as the CLI (as used in the OCI Default Connection):

```hcl
# default
connection "oci" {
  plugin      = "oci"
}
````
