---
title: "Steampipe Table: oci_bds_bds_instance - Query OCI Big Data Service Instances using SQL"
description: "Allows users to query Big Data Service Instances in Oracle Cloud Infrastructure."
---

# Table: oci_bds_bds_instance - Query OCI Big Data Service Instances using SQL

Oracle Cloud Infrastructure Big Data Service is a fully managed service that allows you to process large amounts of data at the speed of business. It provides a high-performance, secure, and easy-to-use environment for running big data workloads, including machine learning and artificial intelligence, real-time analytics, and data exploration. The service is fully integrated with Oracle Cloud Infrastructure, providing seamless access to its data and services.

## Table Usage Guide

The `oci_bds_bds_instance` table provides insights into Big Data Service Instances within Oracle Cloud Infrastructure. As a data engineer or data scientist, explore instance-specific details through this table, including the state of the instance, the time it was created, and the shape of nodes. Utilize it to uncover information about instances, such as their network configuration, whether they are in high availability mode, and the type of workload they are designed for.

## Examples

### Basic info
Explore the configuration details of your Big Data Service instances in Oracle Cloud Infrastructure. This query helps in understanding the current state, security settings, and other crucial details of each instance for better management and decision-making.

```sql
select
  id,
  display_name,
  is_high_availability,
  is_secure,
  is_cloud_sql_configured,
  nodes,
  number_of_nodes,
  cluster_version,
  network_config,
  cluster_details,
  cloud_sql_details,
  created_by,
  bootstrap_script_url,
  kms_key_id,
  cluster_profile,
  lifecycle_state as state
from
  oci_bds_bds_instance;
```

### Count the number of nodes per instance
Assess the distribution of nodes across various instances to manage resources more efficiently. This can be particularly useful in optimizing workload distribution and identifying potential bottlenecks.

```sql
select
  id,
  display_name,
  number_of_nodes
from
  oci_bds_bds_instance;
```

### List secure clusters
Analyze the settings to understand which clusters are secure in your Oracle Cloud Infrastructure Big Data Service. This can be particularly useful for ensuring compliance with security policies and identifying potential vulnerabilities.

```sql
select
  id,
  display_name,
  lifecycle_state,
  is_secure,
  cluster_version,
  created_by
from
  oci_bds_bds_instance
where
  is_secure;
```

### List highly available clusters
Discover the segments that consist of highly available clusters to ensure your data is always accessible and resilient to failures. This is especially useful in maintaining business continuity and minimizing downtime.

```sql
select
  id,
  display_name,
  cluster_profile,
  time_created,
  lifecycle_state,
  is_high_availability,
  created_by
from
  oci_bds_bds_instance
where
  is_high_availability;
```

### List clusters that have cloud SQL configured
Determine the areas in which clusters have been configured with Cloud SQL to gain insights into the high availability, security, and lifecycle state of these instances. This can help in assessing the elements within your Oracle Cloud Infrastructure Big Data Service for optimal resource management.

```sql
select
  id,
  cluster_profile,
  lifecycle_state,
  is_high_availability,
  is_secure,
  is_cloud_sql_configured
from
  oci_bds_bds_instance
where
  is_cloud_sql_configured;
```

### List cloud SQL details of each cluster
Determine the configuration and settings of each cloud SQL cluster, including details such as shape, IP address, block volume size, and Kerberos mapping. This enables efficient management and monitoring of your cloud SQL clusters.

```sql
select
  id,
  display_name
  cloud_sql_details ->> 'Shape' as shape,
  cloud_sql_details ->> 'ShIpAddressape' as ip_address,
  cloud_sql_details ->> 'BlockVolumeSizeInGBs' as block_volume_size_in_gbs,
  cloud_sql_details ->> 'IsKerberosMappedToDatabaseUsers' as is_kerberos_mapped_to_database_users,
  cloud_sql_details ->> 'KerberosDetails' as kerberos_details
from
  oci_bds_bds_instance;
```

### List network config details of each cluster
Analyze the network configuration of each cluster to understand whether a NAT gateway is required and to pinpoint the specific CIDR block being used. This can be beneficial for assessing network requirements and planning infrastructure changes.

```sql
select
  id,
  network_config ->> 'IsNatGatewayRequired' as nat_gateway_required,
  network_config ->> 'CidrBlock' as cidr_block,
from
  oci_bds_bds_instance;
```

### List node information of each cluster
Determine the characteristics of each node within your clusters to better understand their configuration and performance. This can be useful for troubleshooting or optimizing your cluster's performance and resource allocation.

```sql
select
  id,
  display_name,
  n ->> 'InstanceId' as node_instance_id,
  n ->> 'DisplayName' as node_display_name,
  n ->> 'LifecycleState' as node_lifecycle_state,
  n ->> 'NodeType' as node_type,
  n ->> 'Shape' as node_shape,
  n ->> 'SubnetId' as node_subnet_id,
  n ->> 'IpAddress' as node_ip_address,
  n ->> 'SshFingerprint' as node_ssh_fingerprint,
  n ->> 'AvailabilityDomain' as node_availability_domain,
  n ->> 'FaultDomain' as node_fault_domain,
  n ->> 'TimeCreated' as node_time_created,
  n ->> 'AttachedBlockVolumes' as node_attached_block_volumes,
  n ->> 'Hostname' as node_hostname,
  n ->> 'ImageId' as node_image_id,
  n ->> 'TimeUpdated' as node_time_updated,
  n ->> 'Ocpus' as node_ocpus,
  n ->> 'MemoryInGBs' as node_memory_in_gbs,
  n ->> 'Nvmes' as node_nvmes,
  n ->> 'LocalDisksTotalSizeInGBs' as node_local_disks_total_size_in_gbs
from
  oci_bds_bds_instance,
  jsonb_array_elements(nodes) as n;
```

### List KMS key details of each cluster
Explore which clusters are associated with specific KMS keys to understand their protection mode and management details. This is useful for assessing the security configuration of each cluster.

```sql
select
  i.display_name,
  i.kms_key_id,
  k.vault_id,
  k.management_endpoint,
  k.algorithm,
  k.current_key_version,
  k.protection_mode
from
  oci_bds_bds_instance as i,
  oci_kms_key as k
where
  i.kms_key_id = k.id;
```