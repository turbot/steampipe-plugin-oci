# Table: oci_core_image_custom

An image is a template of a virtual hard drive. The image determines the operating system and other software for an instance.You can create a custom image of a bare metal instance's boot disk and use it to launch other instances. Instances you launch from your image include the customizations, configuration, and software installed when you created the image.

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
  oci_core_image_custom;
```

### List images with encryption in transit disabled

```sql
select
  display_name,
  id,
  launch_options ->> 'isPvEncryptionInTransitEnabled' as is_encryption_in_transit_enabled
from
  oci_core_image_custom
where
  launch_options ->> 'isPvEncryptionInTransitEnabled' = 'false';
```
