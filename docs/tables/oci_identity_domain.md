---
title: "Steampipe Table: oci_identity_domain - Query OCI Identity Service Domain using SQL"
description: "Allows users to query OCI Identity Service Domains."
---

# Table: oci_identity_domain - Query OCI Identity Service Groups using SQL

Oracle Cloud Infrastructure (OCI) Identity and Access Management (IAM) service lets you control who has access to your cloud resources. An identity domain is used to manage users and groups, integration standards, external identities, and secure application integration through Oracle Single Sign-on (SSO) configuration.

## Table Usage Guide

The `oci_identity_domain` table provides insights into the groups within OCI Identity and Access Management (IAM). As a security analyst, you can explore domain-specific details through this table, including the type of domain, license type, replica regions, and other associated metadata. Use it to discover information about domains.

## Examples

### Basic info
Explore which identity domain have been created in your OCI environment, along with their lifecycle states, to understand their current status and when they were established. This could be useful for auditing purposes or for maintaining an overview of your security settings.

```sql+postgres
select
  display_name,
  id,
  description,
  lifecycle_state,
  time_created
from
  oci_identity_domain;
```

```sql+sqlite
select
  display_name,
  id,
  description,
  lifecycle_state,
  time_created
from
  oci_identity_domain;
```

### List of Identity Domains that are not in Active state
Discover the segments that consist of identity domain not currently in an active state. This is beneficial in identifying and managing inactive domains within your Oracle Cloud Infrastructure.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state
from
  oci_identity_domain
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state
from
  oci_identity_domain
where
  lifecycle_state <> 'ACTIVE';
```

### List of Identity Domains without application tag key
Determine the areas in which identity domains lack an application tag key. This is useful for identifying potential gaps in your tagging strategy, helping to ensure all domains are properly categorized and managed.

```sql+postgres
select
  display_name,
  id
from
  oci_identity_domain
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  display_name,
  id
from
  oci_identity_domain
where
  json_extract(tags, '$.application') is null;
```

### Get replication details for the domains
Discover the domains that have multiple replicas. This query can be used to get an overview of the replication status of identity domains, helping in effective management and maintenance, for security assessments and ensuring that replication aligns with the organization's security policies, to verify that data replication meets regulatory requirements, especially in cases where data sovereignty and regional compliance are important, to confirm that replicas are available and in the expected state across designated regions.

```sql+postgres
select
  display_name,
  id,
  r ->> 'region' as replica_region,
  r ->> 'state' as replication_state,
  r ->> 'url' as replication_url
from
  oci_identity_domain,
  jsonb_array_elements(replica_regions) as r;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(r.value, '$.region') as replica_region,
  json_extract(r.value, '$.state') as replication_state,
  json_extract(r.value, '$.url') as replication_url
from
  oci_identity_domain,
  json_each(replica_regions) as r;
```