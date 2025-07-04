---
title: "Steampipe Table: oci_service_catalog_private_application - Query OCI Service Catalog Private Applications using SQL"
description: "Allows users to query OCI Service Catalog Private Applications, providing insights into application configurations, lifecycle states, and metadata."
---

# Table: oci_service_catalog_private_application - Query OCI Service Catalog Private Applications using SQL

Oracle Cloud Infrastructure Service Catalog Private Applications allow users to create and manage private software packages that can be shared within an organization. These private applications can be deployed and managed through the OCI Service Catalog, providing a standardized method for distributing internal applications. The `oci_service_catalog_private_application` table provides insights into private applications within your OCI environment.

## Table Usage Guide

The `oci_service_catalog_private_application` table provides insights into private applications within OCI Service Catalog. As a cloud administrator or application developer, explore application-specific details through this table, including package types, lifecycle states, and associated metadata. Utilize it to monitor and manage private applications, verify their configurations, and ensure they are properly maintained within your organization's cloud infrastructure.

## Examples

### Basic info

Explore the essential details of your private applications in the OCI Service Catalog to gain a better understanding of their current status and configuration. This query helps in identifying applications by name, checking their current state, and understanding what type of packages they contain.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  package_type,
  time_created
from
  oci_service_catalog_private_application;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  package_type,
  time_created
from
  oci_service_catalog_private_application;
```

### List private applications by compartment

Identify private applications organized by their respective compartments to better manage resources and access control. This query helps administrators understand the distribution of private applications across different organizational units.

```sql+postgres
select
  c.name as compartment_name,
  a.display_name as application_name,
  a.id,
  a.lifecycle_state,
  a.package_type
from
  oci_service_catalog_private_application a,
  oci_identity_compartment c
where
  a.compartment_id = c.id
order by
  c.name,
  a.display_name;
```

```sql+sqlite
select
  c.name as compartment_name,
  a.display_name as application_name,
  a.id,
  a.lifecycle_state,
  a.package_type
from
  oci_service_catalog_private_application a
join
  oci_identity_compartment c
on
  a.compartment_id = c.id
order by
  c.name,
  a.display_name;
```

### Find private applications created in the last 30 days

Monitor recently created private applications to track new additions to your Service Catalog. This helps in overseeing recent changes and ensuring newly added applications meet organizational standards.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  package_type,
  time_created
from
  oci_service_catalog_private_application
where
  time_created > now() - interval '30 days'
order by
  time_created desc;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  package_type,
  time_created
from
  oci_service_catalog_private_application
where
  time_created > datetime('now', '-30 days')
order by
  time_created desc;
```

### List private applications with their descriptions

Analyze private applications along with their short and long descriptions to better understand their purpose and functionality. This information is valuable for users browsing the Service Catalog to find appropriate applications.

```sql+postgres
select
  display_name,
  id,
  short_description,
  long_description,
  package_type
from
  oci_service_catalog_private_application;
```

```sql+sqlite
select
  display_name,
  id,
  short_description,
  long_description,
  package_type
from
  oci_service_catalog_private_application;
```

### Get private applications with specific tags

Identify private applications that have been tagged with specific metadata to support organizational categorization and governance. This query helps in filtering applications based on business-specific tag classifications.

```sql+postgres
select
  display_name,
  id,
  tags,
  package_type
from
  oci_service_catalog_private_application
where
  tags ? 'Environment'
  and tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  display_name,
  id,
  tags,
  package_type
from
  oci_service_catalog_private_application
where
  json_extract(tags, '$.Environment') = 'Production';
```
