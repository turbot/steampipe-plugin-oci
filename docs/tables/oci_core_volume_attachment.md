---
title: "Steampipe Table: oci_core_volume_attachment - Query OCI Core Volume Attachments using SQL"
description: "Allows users to query OCI Core Volume Attachments."
---

# Table: oci_core_volume_attachment - Query OCI Core Volume Attachments using SQL

A Volume Attachment in Oracle Cloud Infrastructure (OCI) represents the relationship between a volume and an instance. It is a component of OCI Core Services that provides block storage capacity for instances. The volume can be attached in read/write or read-only access mode.

## Table Usage Guide

The `oci_core_volume_attachment` table provides insights into volume attachments within OCI Core Services. As a cloud engineer, you can explore details of volume attachments through this table, including the volume and instance they are attached to, the attachment type, and the state of the attachment. Utilize it to manage and monitor the block storage capacity of your instances, and to ensure optimal configuration and usage of your volumes.

## Examples

### Basic info
Explore which volume attachments in your Oracle Cloud Infrastructure are active or inactive by checking their lifecycle states and when they were created. This can help you manage your resources more effectively by identifying unused or outdated attachments.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment;
```

### List idle volume attachments
Determine the areas in which volume attachments are not actively being used. This can help manage resources more effectively by identifying unused elements.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment
where
  lifecycle_state <> 'ATTACHED';
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment
where
  lifecycle_state <> 'ATTACHED';
```

### List read only volume attachments
Explore which volume attachments are set to read-only status. This can be useful to ensure data integrity by preventing unauthorized modifications.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment
where
  is_read_only;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment
where
  is_read_only = 1;
```