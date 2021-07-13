# Table: oci_objectstorage_bucket

The Oracle Cloud Infrastructure Object Storage service is an internet-scale, high-performance storage platform that offers reliable and cost-efficient data durability

## Examples

### Basic info

```sql
select
  name,
  id,
  namespace,
  storage_tier,
  is_read_only
from
  oci_objectstorage_bucket;
```


### List public buckets

```sql
select
  id,
  name,
  namespace,
  public_access_type
from
  oci_objectstorage_bucket
where
  public_access_type LIKE 'Object%';
```


### List buckets with versioning disabled

```sql
select
  id,
  name,
  namespace,
  versioning
from
  oci_objectstorage_bucket
where
  versioning = 'Disabled';
```


### List buckets with object events disabled

```sql
select
  id,
  name,
  namespace,
  object_events_enabled
from
  oci_objectstorage_bucket
where
  not object_events_enabled;
```


### List buckets with replication disabled

```sql
select
  id,
  name,
  namespace,
  replication_enabled
from
  oci_objectstorage_bucket
where
  not replication_enabled;
```

### List buckets without lifecycle

```sql
select
  name,
  id,
  object_lifecycle_policy -> 'items' as object_lifecycle_policy_rules
from
  oci_objectstorage_bucket
where
  object_lifecycle_policy ->> 'items' is null
  or jsonb_array_length(object_lifecycle_policy -> 'items') = 0;
```
