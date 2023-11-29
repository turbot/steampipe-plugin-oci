---
title: "Steampipe Table: oci_core_volume_group - Query OCI Core Volume Groups using SQL"
description: "Allows users to query information about OCI Core Volume Groups."
---

# Table: oci_core_volume_group - Query OCI Core Volume Groups using SQL

An OCI Core Volume Group is a resource within Oracle Cloud Infrastructure that allows users to manage block storage volumes as a single entity. You can use volume groups to batch manage the lifecycle of your block volumes, including backups, clones, and volume group exports. Volume Groups provide a simple and scalable solution to manage and monitor block storage volumes.

## Table Usage Guide

The `oci_core_volume_group` table provides insights into the volume groups within Oracle Cloud Infrastructure's Core service. As a Cloud Engineer, explore volume group-specific details through this table, including volume group state, availability domain, and associated metadata. Utilize it to uncover information about volume groups, such as their size, the volumes they contain, and their backup policies.

## Examples

### Basic info
Explore which volume groups are currently active by assessing their lifecycle state and creation times. This can help manage resources more effectively by identifying older, potentially unused groups.

```sql
select
  id as volume_group_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_group;
```

### List volume groups in a faulty state
Identify instances where volume groups are in a faulty state to quickly address potential issues and maintain system efficiency.

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_group
where
  lifecycle_state = 'FAULTY';
```

### List volume groups with size greater than 1024 GB
Explore which volume groups exceed a size of 1024 GB to manage storage allocation effectively and prevent potential capacity issues.

```sql
select
  id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_volume_group
where
  size_in_gbs > 1024;
```

### List volume groups created in the root compartment
Explore which volume groups within the root compartment are created. This can be useful for understanding the distribution and organization of your data storage.

```sql
select
  id,
  display_name,
  lifecycle_state,
  tenant_id,
  compartment_id
from
  oci_core_volume_group
where
  compartment_id = tenant_id;
```

### List volume groups created in the last 30 days
Explore the recently created volume groups in the last month to manage resources more effectively and to keep track of changes in your infrastructure. This can be particularly useful for monitoring resource allocation and infrastructure scaling over time.

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created,
  size_in_gbs
from
  oci_core_volume_group
where
  time_created >= now() - interval '30' day;
```

### Get volume details for the volume groups
Explore the specifics of volume groups and their associated volumes to better manage storage resources. This can be particularly useful in understanding the auto-tuned VPUs per GB for each volume within a group, aiding in resource optimization.

```sql
select
  g.id as volume_group_id,
  g.display_name as volume_group_diplay_name,
  v.id as volume_id,
  v.auto_tuned_vpus_per_gb
from
  oci_core_volume_group as g,
  oci_core_volume as v,
  jsonb_array_elements_text(volume_ids) as i
where
  v.id = i;
```