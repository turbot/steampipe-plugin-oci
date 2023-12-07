---
title: "Steampipe Table: oci_bastion_session - Query OCI Bastion Sessions using SQL"
description: "Allows users to query OCI Bastion Sessions."
---

# Table: oci_bastion_session - Query OCI Bastion Sessions using SQL

Oracle Cloud Infrastructure (OCI) Bastion Service provides secure, controlled access to target resources located inside private networks. It is a managed SSH session service that provides a secure way to access hosts located inside your virtual cloud network (VCN). The OCI Bastion Service is designed to provide a secure and scalable method for users to access their infrastructure without exposing it to the public internet.

## Table Usage Guide

The `oci_bastion_session` table provides insights into Bastion Sessions within OCI. As a Security or System Administrator, explore session-specific details through this table, including session type, target resource details, and session status. Utilize it to monitor and manage the secure access to your infrastructure, ensuring the sessions are properly managed and no unauthorized access is happening.

## Examples

### Basic info
Explore the details of active sessions in your cloud bastion service. This query helps you understand the configuration of each session, including session duration and associated user, providing a comprehensive view for better management and security oversight.

```sql+postgres
select
  s.id,
  s.bastion_id,
  s.display_name,
  s.bastion_name,
  s.target_resource_details,
  s.key_details,
  s.session_ttl_in_seconds,
  s.bastion_user_name,
  s.ssh_metadata,
  s.key_type,
  s.lifecycle_state as state
from
  oci_bastion_session s
  inner join
    oci_bastion_bastion b
    on b.id = s.bastion_id;
```

```sql+sqlite
select
  s.id,
  s.bastion_id,
  s.display_name,
  s.bastion_name,
  s.target_resource_details,
  s.key_details,
  s.session_ttl_in_seconds,
  s.bastion_user_name,
  s.ssh_metadata,
  s.key_type,
  s.lifecycle_state as state
from
  oci_bastion_session s
  inner join
    oci_bastion_bastion b
    on b.id = s.bastion_id;
```

### Show port forwarding bastion sessions
Explore the details of your active port forwarding sessions through bastion. This is particularly useful for maintaining secure connections and managing SSH access to your resources.

```sql+postgres
select
  id,
  bastion_id,
  display_name,
  bastion_name,
  target_resource_details,
  key_details,
  session_ttl_in_seconds,
  bastion_user_name,
  ssh_metadata,
  key_type,
  lifecycle_state as state
from
  oci_bastion_session
where
  bastion_id = 'ocid'
  and target_resource_details -> 'sessionType' = '"MANAGED_SSH"';
```

```sql+sqlite
select
  id,
  bastion_id,
  display_name,
  bastion_name,
  target_resource_details,
  key_details,
  session_ttl_in_seconds,
  bastion_user_name,
  ssh_metadata,
  key_type,
  lifecycle_state as state
from
  oci_bastion_session
where
  bastion_id = 'ocid'
  and json_extract(target_resource_details, '$.sessionType') = 'MANAGED_SSH';
```

### List bastion sessions which are not active
Identify inactive sessions within your bastion setup to manage resources and maintain optimal system performance. This is particularly useful for troubleshooting and ensuring the efficient use of resources.

```sql+postgres
select
  display_name,
  id,
  bastion_name,
  bastion_id,
  time_created,
  lifecycle_state as state
from
  oci_bastion_session
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  display_name,
  id,
  bastion_name,
  bastion_id,
  time_created,
  lifecycle_state as state
from
  oci_bastion_session
where
  lifecycle_state <> 'ACTIVE';
```