# Table: oci_identity_network_source

A group of IP addresses that are allowed to access resources in your tenancy.

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_identity_network_source;
```


### List inactive network sources

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_identity_network_source
where
  lifecycle_state = 'INACTIVE';
```


### List network sources which support public ip address

```sql
select
  name,
  id,
  public_source_list
from
  oci_identity_network_source
where
  jsonb_array_length(public_source_list) > 0;
```


### Get list of allowed VCN OCID and IP range pairs for network sources

```sql
select
  name,
  id,
  vsl ->> 'ipRanges' as ip_ranges,
  vsl ->> 'vcnId' as vcn_id
from
  oci_identity_network_source,
  jsonb_array_elements(virtual_source_list) as vsl;
```