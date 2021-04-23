# Table: oci_dns_zone

A DNS zone is where DNS records for a specific DNS domain are kept.

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state,
  self,
  version,
  serial
from
  oci_dns_zone;
```


### List global scoped DNS zones

```sql
select
  name,
  id,
  scope
from
  oci_dns_zone
where
  scope = 'GLOBAL';
```


### List DNS zones which are not protected

```sql
select
  name,
  id,
  is_protected
from
  oci_dns_zone
where
  not is_protected;
```


### List primary DNS zones

```sql
select
  name,
  id,
  zone_type
from
  oci_dns_zone
where
  zone_type = 'PRIMARY';
```
