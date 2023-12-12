---
title: "Steampipe Table: oci_objectstorage_object - Query OCI Object Storage Objects using SQL"
description: "Allows users to query OCI Object Storage Objects."
---

# Table: oci_objectstorage_object - Query OCI Object Storage Objects using SQL

Oracle Cloud Infrastructure's Object Storage is a scalable, secure, and durable solution for storing large amounts of unstructured data. It provides reliable and cost-efficient data durability and allows users to access data from anywhere on the web. This service is ideal for data backup, archival, and big data analytics.

## Table Usage Guide

The `oci_objectstorage_object` table provides insights into the objects within OCI Object Storage. As a data analyst, you can explore object-specific details through this table, including the object's metadata, storage tier, and associated bucket. Utilize it to uncover information about objects, such as their size, time created, and the last time they were modified.

## Examples

### Basic info
Explore which objects within your Oracle Cloud Infrastructure (OCI) Object Storage have been modified recently. This can help in tracking changes, managing storage space, and maintaining data security.

```sql+postgres
select
  name,
  bucket_name,
  namespace,
  region,
  size,
  md5,
  time_created,
  time_modified
from
  oci_objectstorage_object;
```

```sql+sqlite
select
  name,
  bucket_name,
  namespace,
  region,
  size,
  md5,
  time_created,
  time_modified
from
  oci_objectstorage_object;
```


### List archived objects
Explore which objects in your cloud storage have been archived. This is useful for understanding your data retention and could aid in cost management by identifying potentially unnecessary storage.

```sql+postgres
select
  name,
  bucket_name,
  namespace,
  region,
  archival_state
from
  oci_objectstorage_object
where
  archival_state = 'Archived';
```

```sql+sqlite
select
  name,
  bucket_name,
  namespace,
  region,
  archival_state
from
  oci_objectstorage_object
where
  archival_state = 'Archived';
```