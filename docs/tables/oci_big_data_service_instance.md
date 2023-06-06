# Table: oci_big_data_service_instance

Oracle Big Data Service is a fully managed, automated cloud service that provides enterprises with a cost-effective Hadoop environment. Customers easily create secure and scalable Hadoop-based data lakes that can quickly process large amounts of data.

## Examples

### Basic info

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
  oci_big_data_service_instance;
```

### Count the number of nodes per instance

```sql
select
  id,
  display_name,
  number_of_nodes
from
  oci_big_data_service_instance;
```

#### List clusters that should be set up as secure

```sql
select
  id,
  display_name,
  lifecycle_state,
  is_secure,
  cluster_version,
  created_by
from
  oci_big_data_service_instance
where
  is_secure;
```

### List highly available clusters

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
  oci_big_data_service_instance
where
  is_high_availability;
```

### List clusters that has cloud sql configured

```sql
select
  id,
  cluster_profile,
  lifecycle_state,
  is_high_availability,
  is_secure,
  is_cloud_sql_configured
from
  oci_big_data_service_instance
where
  is_cloud_sql_configured;
```

### Get cloud sql details of each cluster

```sql
select
  id,
  display_name
  cloud_sql_details ->> 'Shape' as shape,
  cloud_sql_details ->> 'ShIpAddressape' as ip_address,
  cloud_sql_details ->> 'BlockVolumeSizeInGBs' as block_volume_size_in_g_bs,
  cloud_sql_details ->> 'IsKerberosMappedToDatabaseUsers' as is_kerberos_mapped_to_database_users,
  cloud_sql_details ->> 'KerberosDetails' as kerberos_details
from
  oci_big_data_service_instance;
```

### Get network config details of clusters

```sql
select
  id,
  network_config ->> 'IsNatGatewayRequired' as IsNatGatewayRequired,
  network_config ->> 'CidrBlock' as CidrBlock,
from
  oci_big_data_service_instance;
```

### Get node information of each cluster

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
  n -> 'AttachedBlockVolumes' as node_attached_block_volumes,
  n ->> 'Hostname' as node_hostname,
  n ->> 'ImageId' as node_image_id,
  n ->> 'TimeUpdated' as node_time_updated,
  n ->> 'Ocpus' as node_ocpus,
  n ->> 'MemoryInGBs' as node_memory_in_gbs,
  n ->> 'Nvmes' as node_nvmes,
  n ->> 'LocalDisksTotalSizeInGBs' as node_local_disks_total_size_in_gbs
from
  oci_big_data_service_instance,
  jsonb_array_elements(nodes) as n;
```

### Get KMS key details of each cluster

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
  oci_big_data_service_instance as i,
  oci_kms_key as k
where
  i.kms_key_id = k.id;
```
