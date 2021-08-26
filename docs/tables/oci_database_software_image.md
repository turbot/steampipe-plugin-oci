# Table: oci_database_software_image

Database software images give you the ability to create a customized Oracle Database software configuration that includes your chosen updates (PSU, RU or RUR), and optionally, a list of one-off (or interim) patches or an Oracle Home inventory file.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image;
```

### List deleted database software images

```sql
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image
where
  lifecycle_state = 'DELETED';
```

### List database software images older than 90 days

```sql
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```
