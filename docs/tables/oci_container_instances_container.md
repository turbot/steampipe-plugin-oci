---
title: "Steampipe Table: oci_container_instances_container - Query OCI Container Engine for Kubernetes Container using SQL"
description: "Allows users to query OCI Container Engine for Kubernetes Containers."
---

# Table: oci_container_instances_container - Query OCI Container Engine for Kubernetes Container using SQL

OCI Container Engine for Kubernetes (OKE) is a fully-managed, scalable, and highly available service that you can use to deploy your containerized applications to the cloud. Use OKE when your development team wants to reliably build, deploy, and manage cloud-native applications. You specify the compute resources that your applications require, and OKE provisions and manages the underlying compute instances for you.

## Table Usage Guide

The `oci_container_instances_container` table provides insights into the OCI Container Engine for Kubernetes Containers. As a DevOps engineer, you can use this table to explore details of your containerized applications, including the compute resources they require and how OKE manages these resources. Use this table to monitor the deployment and management of your cloud-native applications and ensure they are running optimally.

## Examples

### Basic info
Explore which container instances are currently active and when they were created to gain insights into your OCI Container Instances' lifecycle and availability. This can help identify instances that may require updates or maintenance.

```sql+postgres
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container;
```

```sql+sqlite
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container;
```

### List failed containers
Uncover the details of containers that have failed. This assists in identifying and addressing issues that may be causing the failure, improving overall system performance and stability.

```sql+postgres
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container
where
  lifecycle_state = 'FAILED';
```

```sql+sqlite
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container
where
  lifecycle_state = 'FAILED';
```

### Get resource config details for containers
Explore the configuration details of your container resources to understand their current lifecycle state and capacity limits in terms of virtual CPUs and memory. This can help in managing and optimizing resource use within your cloud environment.

```sql+postgres
select
  display_name,
  lifecycle_state,
  resource_config ->> 'vcpusLimit' as vcpus_limit,
  resource_config ->> 'memoryLimitInGBs' as memory_limit_in_gbs
from
  oci_container_instances_container;
```

```sql+sqlite
select
  display_name,
  lifecycle_state,
  json_extract(resource_config, '$.vcpusLimit') as vcpus_limit,
  json_extract(resource_config, '$.memoryLimitInGBs') as memory_limit_in_gbs
from
  oci_container_instances_container;
```

### List containers where the resource principal is disabled
Discover the segments that have their resource principal disabled. This is useful for identifying potential security risks and ensuring that all resources are properly managed.

```sql+postgres
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container
where
  is_resource_principal_disabled;
```

```sql+sqlite
select
  display_name,
  id,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container
where
  is_resource_principal_disabled = 1;
```

### Get volume mount details for containers
Explore the configuration of your container's volume mounts to understand their lifecycle state and access permissions. This can help manage storage resources effectively and ensure proper security measures are in place.

```sql+postgres
select
  display_name,
  lifecycle_state,
  vm ->> 'mountPath' as mount_path,
  vm ->> 'volumeName' as volume_name,
  vm ->> 'subPath' as sub_path,
  vm ->> 'isReadOnly' as is_read_only,
  vm ->> 'partition' as partition
from
  oci_container_instances_container,
  jsonb_array_elements(volume_mounts) as vm;
```

```sql+sqlite
select
  display_name,
  lifecycle_state,
  json_extract(vm.value, '$.mountPath') as mount_path,
  json_extract(vm.value, '$.volumeName') as volume_name,
  json_extract(vm.value, '$.subPath') as sub_path,
  json_extract(vm.value, '$.isReadOnly') as is_read_only,
  json_extract(vm.value, '$.partition') as partition
from
  oci_container_instances_container,
  json_each(volume_mounts) as vm;
```

### Get health check details for containers
Explore the health status of your containers to ensure they are functioning optimally. This query is particularly useful for maintaining system performance and identifying any potential issues early.

```sql+postgres
select
  display_name,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url,
  jsonb_pretty(health_checks) as health_checks
from
  oci_container_instances_container;
```

```sql+sqlite
select
  display_name,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url,
  health_checks
from
  oci_container_instances_container;
```

### List containers which does not have any environment variables
Discover the segments that have containers without any environment variables set, which could be crucial for certain applications to function correctly. This query enables you to assess the elements within your infrastructure that might need additional configuration for optimal performance.

```sql+postgres
select
  display_name,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container
where
  environment_variables is null;
```

```sql+sqlite
select
  display_name,
  availability_domain,
  lifecycle_state,
  time_created,
  image_url
from
  oci_container_instances_container
where
  environment_variables is null;
```

### List details of the instance to which the container is connected
Uncover the details of your container's associated instance, such as its name, availability domain, lifecycle state, and creation time. This is useful to understand the operational context and status of your container within the Oracle Cloud Infrastructure (OCI) environment.

```sql+postgres
select
  c.display_name as container_name,
  i.id as instance_id,
  i.display_name as instance_name,
  i.availability_domain as instance_availability_domain,
  i.lifecycle_state as instance_lifecycle_state,
  i.time_created instance_create_time
from
  oci_container_instances_container as c,
  oci_container_instances_container_instance as i
where
  c.container_instance_id = i.id;
```

```sql+sqlite
select
  c.display_name as container_name,
  i.id as instance_id,
  i.display_name as instance_name,
  i.availability_domain as instance_availability_domain,
  i.lifecycle_state as instance_lifecycle_state,
  i.time_created instance_create_time
from
  oci_container_instances_container as c
join
  oci_container_instances_container_instance as i
on
  c.container_instance_id = i.id;
```