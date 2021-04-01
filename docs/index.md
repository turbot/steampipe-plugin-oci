---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/oci.svg"
brand_color: "#F80000"
display_name: "Oracle Cloud Infrastructure"
short_name: "oci"
description: "Steampipe plugin for Oracle Cloud Infrastructure services and resource types."
og_description: Query Oracle Cloud with SQL! Open source CLI. No DB required. 
og_image: "/images/plugins/turbot/oci-social-graphic.png"
---

# Oracle Cloud + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[Oracle Cloud](https://oci.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis. 

For example:

```sql
select
  name,
  id,
  is_mfa_activated,
from
  oci_identity_user;
```

```
+-----------------+------------------------+------------------+
| name            | id                     | is_mfa_activated |
+-----------------+------------------------+------------------+
| pam_beesly      | ocid1.user.oc1.aaaa... | false            |
| creed_bratton   | ocid1.user.oc1.aaaa... | true             |
| stanley_hudson  | ocid1.user.oc1.aaaa... | false            |
| michael_scott   | ocid1.user.oc1.aaaa... | false            |
| dwight_schrute  | ocid1.user.oc1.aaaa... | true             |
+-----------------+------------------------+------------------+
```

## Documentation

- **[Table definitions & examples →](oci/tables)**

## Get started

### Install

Download and install the latest Oracle Cloud plugin:

```bash
steampipe plugin install oci
```

### Credentials

| Item | Description |
| - | - |
| Credentials | Create API keys for your user and add to default OCI configuration: ~/.oci/config |
| Permissions | Use policy builder to enable your group with the permission: `Allow group {group_name} to inspect all-resources in tenancy`  |
| Radius | Each connection represents a single OCI Tenant. |
| Resolution |  1. Static credentials in the configuration file with the `tenancy_ocid`, `user_ocid`, `fingerprint` and `private_key_path arguments`..<br />2. Named profile from an OCI config file(~/.oci/config) with the config_file_profile argument.<br />3. Named profile containing security token.<br />4. Instance Principal based authentication. Note: this configuration will only work when run from an OCI instance.<br />5.  If no credentials are specified, the plugin will use the OCI Default Connection |

### Configuration

Installing the latest oci plugin will create a config file (`~/.steampipe/config/oci.spc`) with a single connection named `oci`:

```hcl
connection "oci_tenant_y" {
  plugin                = "oci"
  config_file_profile   = "DEFAULT"          # Name of the profile 
  config_path           = "~/.oci/config"    # Path to config file
  regions               = ["ap-mumbai-1" , "us-ashburn-1"] # List of regions
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-oci
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)