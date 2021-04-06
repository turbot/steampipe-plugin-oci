# Table: oci_core_internet_gateway

An internet gateway is an optional virtual router you can add to your VCN to enable direct connectivity to the internet.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_internet_gateway;
```


### List disabled internet gateways

```sql
select
  display_name,
  id,
  time_created,
  is_enabled
from
  oci_core_internet_gateway
where
  not is_enabled;
```