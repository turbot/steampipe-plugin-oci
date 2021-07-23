# Table: oci_mysql_db_system_metric_connections_daily

OCI Monitoring metrics explorer provide data about the performance of your systems.  The `oci_mysql_db_system_metric_connections_daily` table provides metric statistics at 24 hour intervals for the most recent 90 days.

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
  oci_mysql_db_system_metric_connections_daily
order by
  id,
  timestamp;
```

### List statistics for ActiveConnections

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_daily
where metric_name = 'ActiveConnections'
order by
  id,
  timestamp;
```

### List statistics for CurrentConnections

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_daily
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```