# Table: oci_mysql_heat_wave_cluster

A HeatWave cluster is a database accelerator for a DB System. A HeatWave Cluster includes a MySQL DB System node and two or more HeatWave nodes. The MySQL DB System node has a HeatWave plugin that is responsible for cluster management, loading data into the HeatWave Cluster, query scheduling, and returning query results to the MySQL DB System.

## Examples

### Basic info

```sql
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster;
```

### List Heat Wave clusters that are not active

```sql
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster
where
  lifecycle_state <> 'ACTIVE';
```

### List heat wave clusters with cluster size more than 2

```sql
select
  db_system_id,
  lifecycle_state as state,
  time_created,
  cluster_size
from
  oci_mysql_heat_wave_cluster
where
  cluster_size > 2;
```
