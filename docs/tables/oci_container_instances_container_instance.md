---
title: "Steampipe Table: oci_container_instances_container_instance - Query OCI Container Engine for Kubernetes Container Instances using SQL"
description: "Allows users to query OCI Container Engine for Kubernetes Container Instances."
---

# Table: oci_container_instances_container_instance - Query OCI Container Engine for Kubernetes Container Instances using SQL

The OCI Container Engine for Kubernetes (OKE) is a fully-managed, scalable, and highly available service that you can use to deploy your containerized applications to the cloud. It provides developers with the ability to build, deploy, and manage cloud-native applications. It is an integral part of the Oracle Cloud Infrastructure platform, providing a seamless experience for managing your cloud infrastructure and applications.

## Table Usage Guide

The `oci_container_instances_container_instance` table provides insights into container instances within OCI Container Engine for Kubernetes (OKE). As a DevOps engineer, you can explore specific details of your container instances through this table, including their metadata, status, and associated resources. Use this table to manage and monitor your containerized applications, ensuring their optimal performance and availability.

## Examples

### Basic info
Explore which container instances are available in your OCI environment. This query can help you assess the lifecycle state and creation time of each instance, and how many containers each instance is running, providing a quick overview of your OCI infrastructure.

```sql
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  container_count
from
  oci_container_instances_container_instance;
```

### List container instances created in the last 7 days
Explore which container instances have been created in the past week. This can be useful for tracking recent activity and managing resource allocation.

```sql
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  container_count
from
  oci_container_instances_container_instance
where
  time_created >= now() - interval '7' day;
```

### List failed container instances
Discover the segments that have failed container instances to assess potential issues and manage system resources more effectively. This can help in troubleshooting and maintaining the overall health of your infrastructure.

```sql
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  container_count
from
  oci_container_instances_container_instance
where
  lifecycle_state = 'FAILED';
```

### Get shape config details for container instances
Explore the configuration details of your container instances to understand their processing capabilities and network bandwidth. This can aid in resource management and optimization of your application's performance.

```sql
select
  display_name,
  lifecycle_state,
  shape_config ->> 'ocpus' as ocpus,
  shape_config ->> 'memoryInGBs' as memory_in_gbs,
  shape_config ->> 'processorDescription' as processor_description,
  shape_config ->> 'networkingBandwidthInGbps' as networking_bandwidth_in_gbps
from
  oci_container_instances_container_instance;
```

### List container instances with no restart policy
Uncover the details of container instances that lack a restart policy, which may be crucial in maintaining service availability and preventing unforeseen downtime. This query is particularly useful for identifying potential vulnerabilities in system resilience.

```sql
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  container_count
from
  oci_container_instances_container_instance
where
  container_restart_policy = 'NEVER';
```

### Get DNS config details for container instances
Explore the DNS configuration details of your container instances to understand their lifecycle state and specific network settings. This can be useful in troubleshooting network issues or optimizing network performance.

```sql
select
  display_name,
  lifecycle_state,
  dns_config -> 'nameservers' as nameservers,
  dns_config -> 'searches' as searches,
  dns_config -> 'options' as options
from
  oci_container_instances_container_instance;
```

### List container instances that are not attached to any volume
Identify container instances that are not linked to any volume in order to assess potential underutilization or misconfiguration. This can be useful for managing resources and optimizing your cloud infrastructure.

```sql
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  container_count
from
  oci_container_instances_container_instance
where
  volume_count is null;
```

### List container instances having more than one containers associated to it
Determine the areas in which a container instance has more than one associated container. This is useful for managing resources and troubleshooting potential issues in configurations.

```sql
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  container_count
from
  oci_container_instances_container_instance
where
  container_count > 1;
```