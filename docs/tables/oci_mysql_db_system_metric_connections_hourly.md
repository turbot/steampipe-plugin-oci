# Table: oci_mysql_db_system_metric_connections_hourly

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_mysql_db_system_metric_connections_hourly` table provides metric statistics at 60 minutes intervals for the most recent 60 days.

## Examples

### Basic info

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_connections_hourly
order by
  id,
  timestamp;
```

### Active connection statistics

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_hourly
where metric_name = 'ActiveConnections'
order by
  id,
  timestamp;
```

### Current connection statistics

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_hourly
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```
