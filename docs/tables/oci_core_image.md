# Table: oci_core_image

An image is a template of a virtual hard drive. The image determines the operating system and other software for an instance.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  size_in_mbs,
  tags,
  lifecycle_state,
  operating_system
from
  oci_core_image;
```

### List images with encryption in transit disabled

```sql
select
  display_name,
  id,
  launch_options ->> 'isPvEncryptionInTransitEnabled' as is_encryption_in_transit_enabled
from
  oci_core_image
where
  launch_options ->> 'isPvEncryptionInTransitEnabled'  = 'false';
```
