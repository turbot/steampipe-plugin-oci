---
title: "Steampipe Table: oci_core_instance - Query OCI Core Instances using SQL"
description: "Allows users to query OCI Core Instances."
---

# Table: oci_core_instance - Query OCI Core Instances using SQL

Oracle Cloud Infrastructure (OCI) Core Instances are virtual server hosts in the cloud. They are the basic compute unit and can be customized based on the needs of the application. Core Instances help in running applications, hosting websites, and supporting large scale analytics scenarios.

## Table Usage Guide

The `oci_core_instance` table provides insights into the Core Instances within Oracle Cloud Infrastructure (OCI). As a Cloud Engineer, you can explore instance-specific details through this table, including instance configurations, state, and associated metadata. Utilize it to uncover information about instances, such as those with specific shapes, the availability domains of instances, and the verification of instance configurations.

## Examples

### Basic info
Discover the segments that are crucial in understanding the status and location of instances in Oracle Cloud Infrastructure. This can be useful in managing resources and optimizing performance across different regions.

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  shape,
  region
from
  oci_core_instance;
```

### List instances along with the compartment details
Determine the areas in which specific instances are located and their associated compartment details, helping to better manage and organize your resources.

```sql
select
  inst.display_name,
  inst.id,
  inst.shape,
  inst.region,
  comp.name as compartment_name
from
  oci_core_instance inst
  inner join
    oci_identity_compartment comp
    on (inst.compartment_id = comp.id)
order by
  comp.name,
  inst.region,
  inst.shape;
```

### Count the number of instances by shape
Analyze the distribution of instances based on their shape to understand the usage pattern and optimize resource allocation. This can help in making informed decisions about infrastructure scaling and cost management.

```sql
select
  shape,
  count(shape) as count
from
  oci_core_instance
group by
  shape;
```

### List instances with shape configuration details
Explore instances to understand their configuration details such as maximum VNIC attachments, memory, networking bandwidth, OCPUs, GPU, and local disk size. This can help in assessing the resources consumed by these instances and planning for future resource allocation.

```sql
select
  display_name,
  shape_config_max_vnic_attachments,
  shape_config_memory_in_gbs,
  shape_config_networking_bandwidth_in_gbps,
  shape_config_ocpus,
  shape_config_baseline_ocpu_utilization,
  shape_config_gpus,
  shape_config_local_disks,
  shape_config_local_disks_total_size_in_gbs
from
  oci_core_instance;
```