# Table: oci_core_vcn

A virtual private network that you set up in Oracle data centers. It closely resembles a traditional network, with firewall rules and specific types of communication gateways that you can choose to use. A VCN resides in a single Oracle Cloud Infrastructure region and covers one or more CIDR blocks of your choice.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state,
  cidr_block,
  freeform_tags
from
  oci_core_vcn;
```

### List unavailable virtual cloud networks

```sql
select
  display_name,
  id,
  lifecycle_state
from
  oci_core_vcn
where
  lifecycle_state <> 'AVAILABLE';
```
