# Table: oci_identity_network_source

A network source is a set of defined IP addresses. The IP addresses can be public IP addresses or IP addresses from VCNs within your tenancy. After you create the network source, you can reference it in policy or in your tenancy's authentication settings to control access based on the originating IP address.

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


### List network sources that include public IP addresses

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


### Get allowed VCN OCIDs and IP range pairs for each network source

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
