# Table: oci_core_image

An image is a template of a virtual hard drive. The image determines the operating system and other software for an instance.

## Examples

### List all images

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


### Get a specific image

```sql
select
  display_name,
  id,
  operating_system,
  operating_system_version,
  freeform_tags
from
  oci_core_image
where
  id = 'ocid1.image.oc1.ap-mumbai-1.aaaaaaaaggwqe5ivg7iayjfko4s3hscukycvvtcsb2gvu2ggeyz7hr3eb1st';
```


### List of images where EncryptionInTransit is not enable

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

### List all custom images

```sql
select
  display_name,
  id,
  lifecycle_state
from
  oci_core_image
where
  tenant_id is not null;
```