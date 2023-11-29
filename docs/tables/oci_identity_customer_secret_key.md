---
title: "Steampipe Table: oci_identity_customer_secret_key - Query OCI Identity Customer Secret Keys using SQL"
description: "Allows users to query OCI Identity Customer Secret Keys."
---

# Table: oci_identity_customer_secret_key - Query OCI Identity Customer Secret Keys using SQL

The Oracle Cloud Infrastructure (OCI) Identity service includes resources that help with the management of identity and access control. One such resource is a Customer Secret Key, which is used for Amazon S3 compatible APIs. These keys enable users to interact with object storage services using tools and libraries that are compatible with S3.

## Table Usage Guide

The `oci_identity_customer_secret_key` table offers insights into the customer secret keys within OCI Identity. As a security analyst, you can explore details about these keys through this table, such as their access keys, associated user IDs, and states. This can be particularly useful for auditing purposes, allowing you to track key usage, identify keys associated with specific users, and monitor the lifecycle states of keys.

## Examples

### Basic info
Explore which customer secret keys have been created in your Oracle Cloud Infrastructure account and when they were created. This can help you manage access to your resources and track account activity over time.

```sql
select
  id,
  display_name,
  user_id,
  user_name,
  time_created
from
  oci_identity_customer_secret_key;
```


### List inactive customer secret keys
Discover the segments that contain inactive customer secret keys to manage access control and enhance security measures. This helps in identifying potential security risks and taking appropriate action.

```sql
select
  id,
  display_name,
  user_id,
  user_name,
  lifecycle_state,
  time_created
from
  oci_identity_customer_secret_key
where
  lifecycle_state = 'INACTIVE';
```


### Count customer secret keys by user
Gain insights into how many secret keys each user has. This query is useful for security audits, ensuring users don't have an excessive number of keys which could increase potential security risks.

```sql
select
  user_id,
  count (id) as customer_secret_key_count
from
  oci_identity_customer_secret_key
group by
  user_id;
```