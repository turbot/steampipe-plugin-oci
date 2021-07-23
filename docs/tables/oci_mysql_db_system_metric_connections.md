# Table: oci_mysql_db_system_metric_connections

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_mysql_db_system_metric_connections` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  oci_mysql_db_system_metric_connections
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
  oci_mysql_db_system_metric_connections
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
  oci_mysql_db_system_metric_connections
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```
