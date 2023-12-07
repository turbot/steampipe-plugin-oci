---
title: "Steampipe Table: oci_core_boot_volume_attachment - Query OCI Core Boot Volume Attachments using SQL"
description: "Allows users to query OCI Core Boot Volume Attachments."
---

# Table: oci_core_boot_volume_attachment - Query OCI Core Boot Volume Attachments using SQL

A Boot Volume Attachment in OCI (Oracle Cloud Infrastructure) is a relationship between a boot volume and an instance, which defines the boot volume from which an instance boots. The boot volume contains the image used to boot the instance, including the operating system and any installed software. The attachment is a critical component for an instance's functionality.

## Table Usage Guide

The `oci_core_boot_volume_attachment` table provides insights into boot volume attachments within OCI Core services. As a Cloud Administrator, explore attachment-specific details through this table, including the state of the attachment, the type of instance attached, and the availability domain. Utilize it to uncover information about boot volume attachments, such as the boot volume's lifecycle state, the time it was created, and the specific OCID (Oracle Cloud Identifier) for the attachment.

## Examples

### Basic info
Discover the segments that are in various stages of their lifecycle and when they were created. This query helps to manage and track resources effectively by providing insights into their current state and creation time.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume_attachment;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume_attachment;
```

### List volume attachments witch are not attached
Analyze the settings to understand the status of volume attachments that are not currently connected. This could be beneficial for managing storage resources and ensuring efficient utilization.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume_attachment
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
  oci_core_boot_volume_attachment
where
  lifecycle_state <> 'ATTACHED';
```