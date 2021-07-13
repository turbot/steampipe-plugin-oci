# Table: oci_core_drg

 A dynamic routing gateway (DRG), which is a virtual router that provides a path for private network traffic between your VCN and existing network.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state,
  time_created
from
  oci_core_drg;
```


### List unavailable dynamic routing gateways

```sql
select
  display_name,
  id,
  lifecycle_state
from
  oci_core_drg
where
  lifecycle_state <> 'AVAILABLE';
```


### Count of dynamic routing gateways per region

```sql
select
  region,
  count(*) drg_count
from
  oci_core_drg
group by
  region;
```