# Table: oci_nosql_table

A NoSQL database service offering on-demand throughput and storage based provisioning that supports JSON, Table and Key-Value datatypes, all with flexible transaction guarantees.

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_nosql_table;
```

### List tables that are not active

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_nosql_table
where
  lifecycle_state <> 'ACTIVE';
```

### List tables with disk storage greater than 1024 GB

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_nosql_table
where
  cast(table_limits -> 'maxStorageInGBs' as INTEGER) > 1024;
```

### Count child tables for parent tables with children

```sql
select
 t2.name as parent,
 count(t1.*) as child_count
from
 oci_nosql_table t1
  join oci_nosql_table t2
  on t1.title like t2.title || '.%'
  and t1.title <> t2.title
group by
 parent;
```

### Count child tables for parent tables with and without children

```sql
select
 t2.name as parent,
 count(t1.*) - 1 as child_count
from
 oci_nosql_table t1
  join oci_nosql_table t2
  on t1.title like t2.title || '%'
group by
 parent;
```