# Table: oci_container_instances_container_instance

Container Instances is a serverless compute service that enables you to quickly and easily run containers without managing any servers. Container Instances runs your containers on serverless compute optimized for container workloads that provides the same isolation as virtual machines. A container instance is defined with a compute shape and one or more containers specifying the container images. Some parameters are defined to configure resource (VCPU, Memory) allocation, networking, and other advanced options such as restart policy and graceful shutdown. Environment variables, startup options, and resource limits can be configured for each container in the container instance.

## Examples

### Basic info

```sql
select
   display_name,
   id,
   time_created,
   containers
from
  oci_container_instances_container_instance;
```

### List active container instances

```sql
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  time_updated,
  shape
from
  oci_container_instances_container_instance
where
  lifecycle_state = 'ACTIVE';
```
### List all available details for container_instances

```sql
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  time_updated,
  time_updated,
  shape,
  containers,
  volume_count,
  graceful_shutdown_timeout_in_seconds,
  fault_domain,
  container_restart_policy,
  availability_domain
  tags
from
  oci_container_instances_container_instance
```