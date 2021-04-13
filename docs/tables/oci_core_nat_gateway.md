# Table: oci_core_nat_gateway

NAT is a networking technique commonly used to give an entire private network access to the internet without assigning each host a public IPv4 address. The hosts can initiate connections to the internet and receive responses, but not receive inbound connections initiated from the internet.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  public_ip_id,
  region,
  tags
from
  oci_core_nat_gateway;
```


### List NAT Gateways that block traffics

```sql
select
  display_name,
  id,
  block_traffic
from
  oci_core_nat_gateway inst
where
  block_traffic;
```


### Count of NAT gateways per VCN

```sql
select
  vcn_id,
  count(*) as nat_gateway_count
from
  oci_core_nat_gateway
group by
  vcn_id;
```
