# Table: oci_core_local_peering_gateway

A local peering gateway (LPG) is an object on a VCN that lets that VCN peer with another VCN in the same region.

## Examples

### Basic info

```sql
select
  name,
  id,
  vcn_id,
  lifecycle_state
from
  oci_core_local_peering_gateway;
```

### List available LPGs

```sql
select
  name,
  id,
  vcn_id,
  lifecycle_state
from
  oci_core_local_peering_gateway
where
  lifecycle_state = 'AVAILABLE';
```

### List LPGs which are not connected to any peer

```sql
select
  name,
  id,
  vcn_id,
  lifecycle_state
from
  oci_core_local_peering_gateway
where
  peering_status = 'NEW';
```
