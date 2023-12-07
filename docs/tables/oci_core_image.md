---
title: "Steampipe Table: oci_core_image - Query OCI Core Images using SQL"
description: "Allows users to query OCI Core Images."
---

# Table: oci_core_image - Query OCI Core Images using SQL

Oracle Cloud Infrastructure (OCI) Core Images are pre-configured operating system images that you can use to create instances in the Compute service. These images include Oracle Linux, CentOS, Ubuntu, and Windows Server, among others. They are designed to provide a secure, stable, and high performance execution environment for applications running on OCI.

## Table Usage Guide

The `oci_core_image` table provides insights into Core Images within Oracle Cloud Infrastructure (OCI). As a system administrator, you can explore image-specific details through this table, including operating system details, launch mode, and associated metadata. Utilize it to uncover information about images, such as their lifecycle state, whether they are bootable, and the time they were created.

## Examples

### Basic info
Explore which images in your Oracle Cloud Infrastructure are currently in use and the operating systems they are running. This information can help you manage your resources more effectively and identify areas for potential optimization or consolidation.

```sql+postgres
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

```sql+sqlite
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

### List images with encryption in transit disabled
Discover the segments that have encryption in transit disabled. This is useful to identify potential security risks and ensure data protection standards are met.

```sql+postgres
select
  display_name,
  id,
  launch_options ->> 'isPvEncryptionInTransitEnabled' as is_encryption_in_transit_enabled
from
  oci_core_image
where
  launch_options ->> 'isPvEncryptionInTransitEnabled' = 'false';
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(launch_options, '$.isPvEncryptionInTransitEnabled') as is_encryption_in_transit_enabled
from
  oci_core_image
where
  json_extract(launch_options, '$.isPvEncryptionInTransitEnabled') = 'false';
```