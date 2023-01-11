# Table: oci_core_cluster_network

A cluster network is a pool of high performance computing (HPC), GPU, or Optimized instances that are connected with a high-bandwidth, ultra low-latency network. Each node in the cluster is a bare metal machine located in close physical proximity to the other nodes. A remote direct memory access (RDMA) network between nodes provides latency as low as single-digit microseconds, comparable to on-premises HPC clusters.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network;
```

### Get instance pool details of cluster network

```sql
select
  c.display_name,
  p -> 'availabilityDomains' as availability_domains,
  p -> 'instanceConfigurationId' as instance_configuration_id,
  p -> 'lifecycleState' as instance_pool_state,
  p -> 'size' as instance_pool_size
from
  oci_core_cluster_network as c,
  jsonb_array_elements(instance_pools) as p;
```

### List available cluster networks

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network
where
  lifecycle_state = 'AVAIALABLE';
```

### Lis cluster networks created in last 30 days

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network
where
  time_created >= now() - interval '30' day;
```