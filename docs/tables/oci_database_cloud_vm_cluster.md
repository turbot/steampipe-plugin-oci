---
title: "Steampipe Table: oci_database_cloud_vm_cluster - Query OCI Database Cloud VM cluster using SQL"
description: "Allows users to query Database Cloud VM cluster within the Oracle Cloud Infrastructure Database service."
---

# Table: oci_database_cloud_vm_cluster - Query OCI Database Cloud VM cluster using SQL

Oracle Cloud Database services, which can be deployed on virtual machines (VMs). A Database Cloud VM Cluster in OCI refers to a specific configuration or deployment option for hosting Oracle databases within a virtualized environment.

## Table Usage Guide

The `oci_database_cloud_vm_cluster` table provides insights into Database Cloud VM cluster within Oracle Cloud Infrastructure Database service. As a database administrator or developer, explore specific details about these VM clusters, including their configurations, statuses, and associated metadata. Utilize this table to manage and monitor your databases, ensuring they are optimized, secure, and functioning as expected.

## Examples

### Basic info
Gain insights into the lifecycle state and creation time of your database cloud VM cluster to better understand their status and duration of existence. This is particularly useful for database management and auditing purposes.

```sql+postgres
select
  cluster_name,
  display_name,
  lifecycle_state,
  time_created,
  shape
from
  oci_database_cloud_vm_cluster;
```

```sql+sqlite
select
  cluster_name,
  display_name,
  lifecycle_state,
  time_created,
  shape
from
  oci_database_cloud_vm_cluster;
```

### List clusters that are not available
Discover the clusters that are currently not available. This can be useful to identify potential issues or disruptions in your cloud database services.

```sql+postgres
select
  cluster_name,
  display_name,
  lifecycle_state,
  time_created,
  availability_domain
from
  oci_database_cloud_vm_cluster
where
  lifecycle_state <> 'AVAILABLE';
```

```sql+sqlite
select
  cluster_name,
  display_name,
  lifecycle_state,
  time_created,
  availability_domain
from
  oci_database_cloud_vm_cluster
where
  lifecycle_state <> 'AVAILABLE';
```

### List VM clusters with core CPUs greater than 4
Identify databases with more than 4 core CPUs. This query is valuable because the number of core CPUs in a computer or server system depends on the specific tasks and processes it executes.

```sql+postgres
select
  cluster_name,
  display_name,
  lifecycle_state,
  time_created,
  availability_domain,
  subnet_id,
  shape,
  cpu_core_count
from
  oci_database_cloud_vm_cluster
where
  cpu_core_count > 4;
```

```sql+sqlite
select
  cluster_name,
  display_name,
  lifecycle_state,
  time_created,
  availability_domain,
  subnet_id,
  shape,
  cpu_core_count
from
  oci_database_cloud_vm_cluster
where
  cpu_core_count > 4;
```

### Get node details of clusters for cloud exadata infrastructures
Queries can assist in making scaling decisions. By analyzing the data from queries, you can determine whether you need to add more nodes to the cluster or remove underutilized nodes to optimize costs and performance.

```sql+postgres
select
  c.cluster_name,
  c.lifecycle_state as cluster_lifecycle_state,
  c.subnet_id,
  i.shape as infra_shape,
  i.cpus_enabled,
  i.storage_count
from
  oci_database_cloud_vm_cluster as c,
  oci_database_exadata_infrastructure as i
where
  c.cloud_exadata_infrastructure_id = i.id;
```

```sql+sqlite
select
  c.cluster_name,
  c.lifecycle_state as cluster_lifecycle_state,
  c.subnet_id,
  i.shape as infra_shape,
  i.cpus_enabled,
  i.storage_count
from
  oci_database_cloud_vm_cluster as c,
  oci_database_exadata_infrastructure as i
where
  c.cloud_exadata_infrastructure_id = i.id;
```