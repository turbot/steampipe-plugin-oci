# Table: oci_container_instances_container_instance

Container Instances is a serverless compute service that enables you to quickly and easily run containers without managing any servers. Container Instances runs your containers on serverless compute optimized for container workloads that provides the same isolation as virtual machines. A container instance is defined with a compute shape and one or more containers specifying the container images. Some parameters are defined to configure resource (VCPU, Memory) allocation, networking, and other advanced options such as restart policy and graceful shutdown. Environment variables, startup options, and resource limits can be configured for each container in the container instance.

## Examples

### Basic info

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

### List failed container instances

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

### Get dns config details for container instances

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