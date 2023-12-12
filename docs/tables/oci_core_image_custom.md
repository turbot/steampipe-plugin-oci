---
title: "Steampipe Table: oci_core_image_custom - Query OCI Core Custom Images using SQL"
description: "Allows users to query custom images in the Oracle Cloud Infrastructure (OCI) Core service."
---

# Table: oci_core_image_custom - Query OCI Core Custom Images using SQL

Custom Images in OCI Core service are user-created images that are used to launch instances within the Oracle Cloud Infrastructure. They are based on existing boot volumes and can be customized to meet specific requirements. These images are critical for creating multiple instances that are consistently configured.

## Table Usage Guide

The `oci_core_image_custom` table provides insights into custom images within the Oracle Cloud Infrastructure Core service. As a DevOps engineer, you can utilize this table to explore details around these images, including their configuration, associated boot volumes, and launch parameters. This table is particularly useful for managing and auditing the custom images used within your OCI environment.

## Examples

### Basic info
Explore the basic details of your custom images in Oracle Cloud Infrastructure, such as display name, size, tags, lifecycle state, and operating system. This can help you manage your resources effectively and keep track of your custom images.

```sql+postgres
select
  display_name,
  id,
  size_in_mbs,
  tags,
  lifecycle_state,
  operating_system
from
  oci_core_image_custom;
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
  oci_core_image_custom;
```

### List images with encryption in transit disabled
Explore which custom images have disabled encryption during data transit, to identify potential security risks. This can help in enhancing the data security measures by ensuring encryption is enabled during data transit.

```sql+postgres
select
  display_name,
  id,
  launch_options ->> 'isPvEncryptionInTransitEnabled' as is_encryption_in_transit_enabled
from
  oci_core_image_custom
where
  launch_options ->> 'isPvEncryptionInTransitEnabled' = 'false';
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(launch_options, '$.isPvEncryptionInTransitEnabled') as is_encryption_in_transit_enabled
from
  oci_core_image_custom
where
  json_extract(launch_options, '$.isPvEncryptionInTransitEnabled') = 'false';
```