# Table: oci_core_dhcp_option

The Networking service provides DHCP options to let you control certain types of configuration on the instances in VCN. Each subnet in a VCN can have a single set of DHCP options associated with it.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  region
from
  oci_core_dhcp_option;
```


### DHCP options configuration parameters info

```sql
select
  id,
  display_name,
  jsonb_array_elements_text(o -> 'searchDomainNames') as search_domain_names,
  jsonb_array_elements_text(o -> 'customDnsServers') as custom_dns_servers,
  o ->> 'serverType' as server_type,
  o ->> 'type' as type
from
  oci_core_dhcp_option,
  jsonb_array_elements(options) as o;
```


### Count of DHCP options per VCN

```sql
select
  vcn_id,
  count(*) dhcp_option_count
from
  oci_core_dhcp_option
group by
  vcn_id;
```
