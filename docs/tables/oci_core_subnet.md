# Table: oci_core_subnet

A logical subdivision of a VCN. Each subnet consists of a contiguous range of IP addresses that do not overlap with other subnets in the VCN.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  tags
from
  oci_core_subnet;
```


### OCIDs of the security list that the subnet uses

```sql
select
  display_name,
  id,
  jsonb_array_elements_text(security_list_ids) as security_list_id
from
  oci_core_subnet;
```


### Subnet count by VCN ID

```sql
select
  vcn_id,
  count(id) as subnet_count
from
  oci_core_subnet
group by
  vcn_id;
```


### Find the number of available IP address in each subnet

```sql
select
  id,
  cidr_block,
  power(2, 32 - masklen(cidr_block :: cidr)) -1 as raw_size
from
  oci_core_subnet;
```
