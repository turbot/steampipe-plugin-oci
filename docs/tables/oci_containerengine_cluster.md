---
title: "Steampipe Table: oci_containerengine_cluster - Query OCI Container Engine Clusters using SQL"
description: "Allows users to query OCI Container Engine Clusters."
---

# Table: oci_containerengine_cluster - Query OCI Container Engine Clusters using SQL

Oracle Container Engine for Kubernetes (OKE) is a fully-managed, scalable, and highly available service that you can use to deploy your containerized applications to the cloud. Use OKE when your development team wants to reliably build, deploy, and manage cloud-native applications. You specify the compute resources that your applications require, and OKE provisions and manages the underlying bare metal instances for you.

## Table Usage Guide

The `oci_containerengine_cluster` table provides insights into Container Engine Clusters within Oracle Cloud Infrastructure (OCI). As a DevOps engineer, explore cluster-specific details through this table, including cluster options, metadata, and associated resources. Utilize it to uncover information about clusters, such as those with specific configurations, the relationships between clusters and nodes, and the verification of cluster policies.

## Examples

### Basic info
Analyze the lifecycle state of your OCI Container Engine clusters to understand their current operational status. This can help in identifying any clusters that may need attention or maintenance.

```sql+postgres
select
  name,
  id,
  lifecycle_state
from
  oci_containerengine_cluster;
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state
from
  oci_containerengine_cluster;
```

### List failed container engine clusters
Explore which container engine clusters have failed to understand potential issues or disruptions in your system. This could assist in troubleshooting, maintenance, and system optimization.

```sql+postgres
select
  name,
  id,
  lifecycle_state
from
  oci_containerengine_cluster
where
  lifecycle_state = 'FAILED';
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state
from
  oci_containerengine_cluster
where
  lifecycle_state = 'FAILED';
```
### List container engine clusters for which image policy is disabled
Explore which container engine clusters have their image policy disabled. This can help identify potential security risks, as disabling image policy may allow unauthorized images to be deployed.

```sql+postgres
select
  name,
  id,
  lifecycle_state,
  image_policy_config_enabled
from
  oci_containerengine_cluster
where
  image_policy_config_enabled = false;
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state,
  image_policy_config_enabled
from
  oci_containerengine_cluster
where
  image_policy_config_enabled = 0;
```