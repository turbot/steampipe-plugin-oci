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

### Get details of specific subnet

```sql
select
  display_name,
  id,
  time_created,
  dns_label,
  cidr_block
from
  oci_core_subnet
where
  id = 'ocid1.subnet.oc1.ap-mumbai-1.aaaaaaaa2rn43msfyjb7k5orwfbypbsuur72xlnv3qybxss5ukherfv5sdfb';
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