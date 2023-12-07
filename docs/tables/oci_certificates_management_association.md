---
title: "Steampipe Table: oci_certificates_management_association - Query OCI Certificates Management Associations using SQL"
description: "Allows users to query OCI Certificates Management Associations."
---

# Table: oci_certificates_management_association - Query OCI Certificates Management Associations using SQL

Oracle Cloud Infrastructure (OCI) Certificates Management Association is a resource within Oracle Cloud that allows users to manage the lifecycle of SSL/TLS certificates. It provides a centralized way to upload, renew, and manage certificates, ensuring the security of your applications and services. OCI Certificates Management Association helps maintain the integrity and confidentiality of your data by facilitating the secure transmission of information over the internet.

## Table Usage Guide

The `oci_certificates_management_association` table provides insights into the associations between certificates and resources within Oracle Cloud Infrastructure (OCI) Certificates Management. As a security administrator, use this table to explore association-specific details, including the certificate's lifecycle state, time of creation, and associated metadata. This table is beneficial for gaining insights into certificate associations, such as those linked with specific resources, the status of these associations, and the verification of certificate expiration dates.

## Examples

### Basic info
Explore the lifecycle state and creation time of your OCI certificates, along with their association types. This can help you manage and track your certificates more effectively.

```sql+postgres
select
  id,
  name,
  lifecycle_state,
  certificates_resource_id,
  association_type,
  time_created
from
  oci_certificates_management_association;
```

```sql+sqlite
select
  id,
  name,
  lifecycle_state,
  certificates_resource_id,
  association_type,
  time_created
from
  oci_certificates_management_association;
```

### Count the number of certificate associations by type
Analyze the number of certificate associations based on their types. This is useful for understanding the distribution and prevalence of different certificate associations within your Oracle Cloud Infrastructure.

```sql+postgres
select
  association_type,
  count(id) as numbers_of_association
from
  oci_certificates_management_association
group by
  association_type;
```

```sql+sqlite
select
  association_type,
  count(id) as numbers_of_association
from
  oci_certificates_management_association
group by
  association_type;
```

### List associations created in the last 10 days
Discover the recently established associations within the last 10 days. This can be beneficial for keeping track of new additions and understanding the recent changes in your system.

```sql+postgres
select
  name,
  id,
  lifecycle_state,
  time_created,
  associated_resource_id
from
  oci_certificates_management_association
where
  time_created >= now() - interval '10' day;
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state,
  time_created,
  associated_resource_id
from
  oci_certificates_management_association
where
  time_created >= datetime('now', '-10 day');
```

### List inactive associations
Explore which certificate management associations are not in an active state. This can be useful to identify and address potential issues or inefficiencies within your system.

```sql+postgres
select
  name,
  id,
  lifecycle_state,
  time_created,
  associated_resource_id
from
  oci_certificates_management_association
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state,
  time_created,
  associated_resource_id
from
  oci_certificates_management_association
where
  lifecycle_state <> 'ACTIVE';
```