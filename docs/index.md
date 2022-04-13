---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/oci.svg"
brand_color: "#F80000"
display_name: "Oracle Cloud Infrastructure"
short_name: "oci"
description: "Steampipe plugin for Oracle Cloud Infrastructure services and resource types."
og_description: "Query Oracle Cloud with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/oci-social-graphic.png"
---

# Oracle Cloud + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[Oracle Cloud](https://www.oracle.com/cloud/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

For example:

```sql
select
  name,
  id,
  is_mfa_activated
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

- **[Table definitions & examples →](/plugins/turbot/oci/tables)**

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
| Permissions | Use policy builder to enable your group with following permissions:<br /><li>`Allow group {group_name} to read all-resources in tenancy`</li><li>`Allow group {group_name} to manage all-resources in tenancy where request.operation='GetConfiguration'`</li>**Note:** Permission to manage `GetConfiguration` for all-resources is required for `oci_identity_tenancy` table. |
| Radius | Each connection represents a single OCI Tenant. |
| Resolution | 1. Static credentials in the configuration file with the `tenancy_ocid`, `user_ocid`, `fingerprint` and `private_key_path arguments`.<br />2. Named profile from an OCI config file(~/.oci/config) with the config_file_profile argument.<br />3. Named profile containing security token.<br />4. Instance Principal based authentication. Note: this configuration will only work when run from an OCI instance.<br />5. If no credentials are specified, the plugin will use the OCI Default Connection |

### Configuration

Installing the latest oci plugin will create a config file (`~/.steampipe/config/oci.spc`) with a single connection named `oci`:

```hcl
connection "oci_tenant_y" {
  plugin              = "oci"

  # Name of the profile.
  #config_file_profile = "DEFAULT"

  # Path to config file
  #config_path = "~/.oci/config"

  # List of regions
  #regions = ["ap-mumbai-1", "us-ashburn-1"]

  # The maximum number of attempts (including the initial call) Steampipe will
  # make for failing API calls. Defaults to 9 and must be greater than or equal to 1.
  #max_error_retry_attempts = 9

  # The minimum retry delay in milliseconds after which retries will be performed.
  # This delay is also used as a base value when calculating the exponential backoff retry times.
  # Defaults to 25ms and must be greater than or equal to 1ms.
  #min_error_retry_delay = 25
}
```

- `config_file_profile` (Optional) OCI profile name to use for credentials.
- `config_path` (Optional) Path of the config file where subjected profile is available.
- `max_error_retry_attempts` (Optional) The maximum number of attempts (including the initial call) Steampipe will make for failing API calls. Defaults to 9 and must be greater than or equal to 1.
- `min_error_retry_delay` (Optional) The minimum retry delay in milliseconds after which retries will be performed. This delay is also used as a base value when calculating the exponential backoff retry times. Defaults to 25ms and must be greater than or equal to 1ms.
- `regions` (Optional) List of OCI regions Steampipe will connect to

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-oci
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)

## Advanced configuration options

If you have an OCI profile setup for using the [OCI CLI](https://docs.oracle.com/en-us/iaas/tools/oci-cli/2.9.1/oci_cli_docs/oci.html), Steampipe will just work with that connection.

For users with multiple accounts and more complex authentication use cases, here are some examples of advanced configuration options:

### Use static credentials

The OCI plugin allows you set static credentials with the tenancy_ocid, user_ocid, fingerprint and private_key_path arguments. You may select one or more regions with the regions argument.

```hcl
connection "oci_tenant_x" {
  plugin           = "oci"
  tenancy_ocid     = "ocid1.tenancy.oc1..aaaaaaaa111111111bbbbbbbetci3yjjnjqmfkr4pab12cd45gh56hm76cyljaq"
  user_ocid        = "ocid1.user.oc1..aaaaaaaa111111111bbbbbbb2oixpabcd7a3jkl6yife75v7a7o6c5d6wclrsjia"
  fingerprint      = "9a:a1:b2:c3:d4:e5:6f:7g:89:33:5f:ed:ab:ec:de:11"
  private_key_path = "~/.ssh/oci_private.pem"          # Path to user's private key
  regions          = ["ap-mumbai-1", "us-ashburn-1"]   # List of regions to query resources
}
```

### Using a named profile

If you have an OCI config file(~/.oci/config) with multiple profiles setup, you can set the config_file_profile argument:

```hcl
connection "oci" {
  plugin                = "oci"
  config_file_profile   = "DEFAULT"          # Name of the profile in the OCI config file
  config_path           = "~/.oci/config"    # Path to config file
  regions               = ["ap-mumbai-1", "us-ashburn-1"] # List of regions to query resources
}

connection "oci_tenant_x" {
  plugin                = "oci"
  config_file_profile   = "tenant_x"         # Name of the profile in the OCI config file
  config_path           = "~/.oci/config"    # Path to config file
  regions               = ["ap-mumbai-1", "us-ashburn-1"] # List of regions to query resources
}
```

### Using a named profile containing security token

```hcl
connection "oci_tenant_z" {
  plugin              = "oci"
  auth                = "SecurityToken"   # Type of authentication
  config_file_profile = "tenant_z"        # OCI Profile containing the details of the token
  regions             = ["ap-mumbai-1"]
}
```

### Instance principal based authentication

This configuration will only work when run from an OCI instance. More information on using [Instance Principals](https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm):

```hcl
connection "oci" {
  plugin = "oci"
  auth   = "InstancePrincipal"   # Type of authentication
}
```
