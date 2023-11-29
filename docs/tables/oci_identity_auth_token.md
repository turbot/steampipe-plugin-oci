---
title: "Steampipe Table: oci_identity_auth_token - Query OCI Identity Auth Tokens using SQL"
description: "Allows users to query OCI Identity Auth Tokens."
---

# Table: oci_identity_auth_token - Query OCI Identity Auth Tokens using SQL

An OCI Identity Auth Token is a feature within Oracle Cloud Infrastructure that allows you to manage and authenticate API requests in OCI services. It provides a secure way to authenticate requests made to OCI resources, including compute instances, databases, and storage services. OCI Identity Auth Tokens help you manage the security and integrity of your OCI resources by providing a means to authenticate requests without exposing your user credentials.

## Table Usage Guide

The `oci_identity_auth_token` table provides insights into Auth Tokens within Oracle Cloud Infrastructure (OCI). As a Security Administrator, explore token-specific details through this table, including its status, description, and associated user details. Utilize it to uncover information about tokens, such as those that are inactive, the users associated with each token, and the lifecycle state of these tokens.

## Examples

### Basic info
Explore which authentication tokens have been created within your Oracle Cloud Infrastructure, along with their associated user details and creation timestamps. This can aid in understanding user activity and tracking token usage.

```sql
select
  id,
  user_id,
  user_name,
  time_created
from
  oci_identity_auth_token;
```


### List inactive auth tokens
Explore which authentication tokens are inactive. This can help in identifying potential security risks, as inactive tokens can be a sign of unauthorized access or outdated user credentials.

```sql
select
  id,
  user_id,
  user_name,
  lifecycle_state,
  time_created
from
  oci_identity_auth_token
where
  lifecycle_state = 'INACTIVE';
```


### Count the number of auth tokens by user
Analyze the settings to understand the distribution of authentication tokens across different users. This is useful to monitor user activity and ensure that no user is generating an excessive number of tokens, which could be a potential security risk.

```sql
select
  user_id,
  user_name,
  count (id) as auth_token_count
from
  oci_identity_auth_token
group by
  user_name,
  user_id;
```


### List auth tokens older than 90 days
Explore which authentication tokens have been active for more than 90 days. This can be useful for identifying potential security risks and maintaining system integrity.

```sql
select
  id,
  user_id,
  user_name,
  lifecycle_state,
  time_created
from
  oci_identity_auth_token
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```