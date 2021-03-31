---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/oci.svg"
brand_color: "#F80000"
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

The plugin needs to be configured with credentials for the Oracle Cloud Infrastructure account.

Default connection configuration:

```hcl
connection "oci" {
  plugin      = "oci"
}
```

Plugin also supports below set of configuration arguments:

1. The OCI plugin allows you set static credentials with the `tenancy_ocid`, `user_ocid`, `fingerprint` and `private_key_path` arguments. You may select one or more regions with the `regions` argument.

   An OCI connection may connect to multiple regions, however be aware that performance may be negatively affected by both the number of regions and the latency to them.

   ```hcl
   # credentials via user API Key pair
   connection "oci_tenant_x" {
     plugin            = "oci"
     tenancy_ocid      = "ocid1.tenancy.oc1..aaaaaaaa111111111bbbbbbbetci3yjjnjqmfkr4pab12cd45gh56hm76cyljaq"
     user_ocid         = "ocid1.user.oc1..aaaaaaaa111111111bbbbbbb2oixpabcd7a3jkl6yife75v7a7o6c5d6wclrsjia"
     fingerprint       = "9a:a1:b2:c3:d4:e5:6f:7g:89:33:5f:ed:ab:ec:de:11"
     private_key_path  = "~/.ssh/oci_private.pem"           # Path to user's private key
     regions           = ["ap-mumbai-1" , "us-ashburn-1"]   # List of regions to query resources
   }
   ```

2. Using a named profile from an OCI config file(`~/.oci/config`) with the `config_file_profile` argument:

   ```hcl
   # credentials via profile
   connection "oci_tenant_y" {
     plugin                = "oci"
     config_file_profile   = "DEFAULT"          # Name of the profile in the OCI config file
     config_path           = "~/.oci/config"    # Path to config file
     regions               = ["ap-mumbai-1" , "us-ashburn-1"] # List of regions to query resources
   }
   ```

3. Using a named profile containing security token

   ```hcl
   connection "oci_tenant_z" {
     plugin              =	"oci"
     auth                =	"SecurityToken"   # Type of authentication.
     config_file_profile =	"profile_with_token" # OCI Profile containing the details of the token
     regions             = ["ap-mumbai-1"]
   }
   ```

4. Configure the Oracle Cloud Infrastructure provider to use Instance Principal based authentication.
   **Note:** this configuration will only work when run from an OCI instance. For more information on using Instance Principals, see this [document](https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm).

   ```hcl
   connection "oci" {
     plugin  =  "oci"
     auth    =  "InstancePrincipal"   # Type of authentication.
   }
   ```

If no credentials are specified, the plugin will use the OCI Default Connection:

```hcl
# default
connection "oci" {
  plugin      = "oci"
}
```

### Order of precedence

The Steampipe OCI plugin respects and applies configurations specified by connection configuration, environment variable, or OCI config file entry in the following order of precedence:

The value specified in the steampipe connection config option.
The value specified in the environment variable.
The value specified in the OCI config file.

If `regions` is not specified, Steampipe will use a single default region using the same resolution order as the credentials:

1. The `OCI_CLI_REGION` or `OCI_REGION` environment variable
2. The region specified in the profile

Steampipe will require read access in order to query your OCI resources.

#### References:

- [Security Credentials](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/credentials.htm)
- [Required IAM Policy to Work with Resources in the Tenancy Explorer](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/compartmentexplorer.htm#iampolicy)
