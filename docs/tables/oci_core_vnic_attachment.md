# Table: oci_core_vnic_attachment

A VNIC attachment represents an attachment between a VNIC and an instance. A VNIC enables an instance to connect to a VCN and determines how the instance connects with endpoints inside and outside the VCN. Each VNIC resides in a subnet in a VCN.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  instance_id,
  availability_domain,
  lifecycle_state,
  private_ip,
  public_ip,
  time_created
from
  oci_core_vnic_attachment;
```

### List detached VNICs

```sql
select
  id,
  display_name,
  lifecycle_state
from
  oci_core_vnic_attachment
where
  lifecycle_state = 'DETACHED';
```

### List automatically created and attached VNICs

```sql
select
  id,
  display_name,
  is_primary
from
  oci_core_vnic_attachment
where
  is_primary;
```

### List VNICs with disabled source/destination check

```sql
select
  id,
  display_name,
  skip_source_dest_check
from
  oci_core_vnic_attachment
where
  skip_source_dest_check;
```
