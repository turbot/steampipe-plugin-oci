# Table: oci_streaming_stream

The Oracle Cloud Infrastructure Streaming service provides a fully managed, scalable, and durable solution for ingesting and consuming high-volume data streams in real-time.

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_streaming_stream;
```

### List streams that are not active

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_streaming_stream
where
  lifecycle_state <> 'ACTIVE';
```

### List streams with retention period greater than 24 hrs

```sql
select
  name,
  id,
  lifecycle_state,
  time_created,
  retention_in_hours
from
  oci_streaming_stream
where
  retention_in_hours > 24;
```
