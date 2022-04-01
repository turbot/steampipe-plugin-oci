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

### List failed heat wave clusters

```sql
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster
where
  lifecycle_state = 'FAILED';
```

### List heat wave clusters older than 90 days

```sql
select
  db_system_id,
  lifecycle_state as state,
  time_created,
  cluster_size
from
  oci_mysql_heat_wave_cluster
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```
