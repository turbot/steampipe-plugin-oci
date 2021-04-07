# Table: oci_core_vcn

A virtual private network that you set up in Oracle data centers. It closely resembles a traditional network, with firewall rules and specific types of communication gateways that you can choose to use. A VCN resides in a single Oracle Cloud Infrastructure region and covers one or more CIDR blocks of your choice.

## Examples

### List all vcn

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

### Get a specific vcn

```sql
select
  display_name,
  id,
  lifecycle_state,
  cidr_block,
  dns_label
from
  oci_core_vcn
where
  id = 'ocid1.vcn.oc1.ap-mumbai-1.amaaaaaa6igdexaasn5aalvhjyqctaiykwy2bg3xjqeza6muotvya1wdrf4v';
```

### List of vcn where state is not available

```sql
select
  display_name,
  id,
  lifecycle_state
from
  oci_core_vcn
where
  lifecycle_state != 'AVAILABLE';
```