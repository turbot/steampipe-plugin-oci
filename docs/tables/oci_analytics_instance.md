# Table: oci_analytics_instance

Oracle Analytics Cloud is a scalable and secure public cloud service that provides a full set of capabilities to explore and perform collaborative analytics for you, your workgroup, and your enterprise. An analytics instance is a virtual server in the Oracle Analytics Cloud. 

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state as state,
  time_created,
  time_updated
from
  oci_analytics_instance;
```

### List active analytics instances

```sql
select
  id,
  name,
  lifecycle_state as state,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  lifecycle_state = 'ACTIVE';
```

### List analytics instances older than 90 days

```sql
select
  id,
  name,
  lifecycle_state,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

### List analytics instances with private network endpoint type

```sql
select
  id,
  name,
  lifecycle_state,
  network_endpoint_details,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  network_endpoint_details ->> 'networkEndpointType' = 'PRIVATE';
```

### Get analytics instance by ID

```sql
select
  id,
  name,
  lifecycle_state,
  network_endpoint_details,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  id = 'ocid1.analyticsinstance.oc1.ap-mumbai-1.aaaaaaaaftab5bh4dp6ng6djpjnhix6mpfj25zwr6z67xnmg4c43k5hu5o2a';
```
