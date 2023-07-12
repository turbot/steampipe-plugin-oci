# Table: oci_logging_search

The service offers powerful search capabilities, enabling you to search for specific log entries based on different criteria, such as time range, log content, log source, or custom filters. This helps you quickly locate relevant logs and investigate issues.

**Important notes:**

- By default, this table will provide data for the last 24hrs. You can give the `timestamp` value in the below ways to fetch data in a range. The examples below can guide you.

  - timestamp >= ‘2023-03-11T00:00:00Z’ and timestamp <= ‘2023-03-15T00:00:00Z’
  - timestamp between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - timestamp > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - timestamp < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)

- We recommend specifying optional quals in the query to optimize the table output. Optional quals are supported for the following columns:
  - `log_group_name`
  - `log_name`
  - `search_query`
  - `timestamp`

## Examples

### Show log entries of last 24 hrs

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search;
```

### Show log entries of last 12 hrs

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  timestamp >= now() - interval '12 hrs';
```

### Show log entries which are related to DatabaseService

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  log_content_type like '%DatabaseService%';
```

### Show log entries from us-ashburn-1 region

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  region = 'us-ashburn-1';
```

### Show log entries from a specific log group

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  log_group_name = 'test-bucket';
```

### Show log entries from a specific search query

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  search_query = 'search "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecx3hoz4p4h2q37cyljaq/test" | sort by datetime desc';
```