# Table: oci_containerengine_cluster

Oracle Cloud Infrastructure Container Engine for Kubernetes is a fully-managed, scalable, and highly available service that you can use to deploy your containerized applications to the cloud

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_containerengine_cluster;
```

### List failed container engine clusters

```sql
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

```sql
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