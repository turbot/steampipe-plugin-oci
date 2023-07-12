# Table: oci_queue_queue

OCI Queue is a fully managed serverless service that helps decouple systems and enable asynchronous operations. OCI Queue service handles high-volume transactional data that requires independently processed messages without loss or duplication. Queue supports transparent, automatic scaling based on throughput for producers and consumers. OCI Queue uses open standards to support communication with any client or producer with minimal effort.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  lifecycle_state
from
  oci_queue_queue;
```

### List queues not in the active state

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_queue_queue
where
  lifecycle_state <> 'ACTIVE';
```

### List queues that are not encrypted

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  compartment_id,
  region
from
  oci_queue_queue
where 
  custom_encryption_key_id is null;
```  

### Get details of each queue

```sql
select
  display_name,
  retention_in_seconds,
  visibility_in_seconds,
  timeout_in_seconds,
  messages_endpoint,
  dead_letter_queue_delivery_count,
  defined_tags,
  id
from
  oci_queue_queue;
```
