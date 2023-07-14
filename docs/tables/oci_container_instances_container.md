# Table: oci_container_instances_container

Containers inside an OCI (Oracle Cloud Infrastructure) Container Instance refer to the individual units of isolated application runtime environments running within the context of the Container Instance. Containers are lightweight and portable, encapsulating all the necessary components for the application to run consistently across different environments.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  display_name,
  lifecycle_state,
  resource_config ->> 'vcpusLimit' as vcpus_limit,
  resource_config ->> 'memoryLimitInGBs' as memory_limit_in_gbs
from
  oci_container_instances_container;
```

### List containers where the resource principal is disabled

```sql
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

### Get volume mount details for containers

```sql
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

### Get health check details for containers

```sql
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

### List containers which does not have any environment variables

```sql
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

```sql
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
