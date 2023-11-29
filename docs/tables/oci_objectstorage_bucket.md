---
title: "Steampipe Table: oci_objectstorage_bucket - Query OCI Object Storage Buckets using SQL"
description: "Allows users to query OCI Object Storage Buckets."
---

# Table: oci_objectstorage_bucket - Query OCI Object Storage Buckets using SQL

Oracle Cloud Infrastructure's Object Storage service is an internet-scale, high-performance storage platform that offers reliable and cost-efficient data durability. The Object Storage service can store an unlimited amount of unstructured data of any content type, including analytic data and rich content, like images and videos. With strong consistency, your data is reliably stored and retrieved.

## Table Usage Guide

The `oci_objectstorage_bucket` table provides insights into Object Storage Buckets within Oracle Cloud Infrastructure's Object Storage service. As a data engineer, you can explore bucket-specific details through this table, including its current state, storage tier, and associated metadata. Utilize it to uncover information about buckets, such as their public accessibility, region, and time of creation.

## Examples

### Basic info
Explore which storage buckets in your cloud environment are set to read-only. This can help you determine areas where data cannot be modified, aiding in data management and security.

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
Explore which storage buckets in your Oracle Cloud Infrastructure have public access. This is useful for identifying potential security risks and ensuring data privacy.

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
Identify the storage buckets where versioning is disabled. This is useful for assessing potential risks, as these buckets don't have the ability to recover previous versions of the data.

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
Determine the areas in which object events are disabled within your data storage. This is useful for identifying potential gaps in your event tracking and monitoring setup.

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
Identify storage buckets where replication is not enabled. This can be useful for ensuring data redundancy and availability in your infrastructure.

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
Discover the segments that lack a lifecycle policy in the object storage buckets. This is useful for identifying and rectifying areas where data might be accumulating indefinitely, leading to unnecessary storage costs.

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