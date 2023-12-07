---
title: "Steampipe Table: oci_bastion_bastion - Query OCI Bastion Service Bastions using SQL"
description: "Allows users to query Bastions in the OCI Bastion Service."
---

# Table: oci_bastion_bastion - Query OCI Bastion Service Bastions using SQL

The OCI Bastion Service provides secure, controlled access to target resources that reside in private networks. It acts as a 'jump host' for administrators to securely access their cloud resources. This service is especially useful for resources that do not have public endpoints.

## Table Usage Guide

The `oci_bastion_bastion` table provides insights into Bastions within the OCI Bastion Service. As a system administrator, you can explore details of each Bastion through this table, including its configuration, status, and associated metadata. Use this table to understand the setup of your Bastions, verify their configurations, and ensure they are providing secure access as expected.

## Examples

### Basic info
Explore the configuration of your bastion host in Oracle Cloud Infrastructure, including its type, status, and associated network details. This can help you manage your cloud security by understanding the maximum sessions allowed, session lifespan, and the state of each bastion host.

```sql+postgres
select
  id,
  name,
  bastion_type,
  dns_proxy_status,
  client_cidr_block_allow_list,
  max_session_ttl_in_seconds,
  max_sessions_allowed,
  private_endpoint_ip_address,
  static_jump_host_ip_address,
  phone_book_entry,
  target_vcn_id,
  target_subnet_id,
  lifecycle_state as state
from
  oci_bastion_bastion;
```

```sql+sqlite
select
  id,
  name,
  bastion_type,
  dns_proxy_status,
  client_cidr_block_allow_list,
  max_session_ttl_in_seconds,
  max_sessions_allowed,
  private_endpoint_ip_address,
  static_jump_host_ip_address,
  phone_book_entry,
  target_vcn_id,
  target_subnet_id,
  lifecycle_state as state
from
  oci_bastion_bastion;
```

### Show Bastions that allow access from the Internet (0.0.0.0/0)
Identify Bastions that permit internet access, providing insights into potential security vulnerabilities within your network infrastructure.

```sql+postgres
select
  id,
  name,
  bastion_type,
  client_cidr_block_allow_list,
  private_endpoint_ip_address
from
  oci_bastion_bastion
where
  (
    client_cidr_block_allow_list
  )
  ::jsonb ? '0.0.0.0/0';
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```

### List bastions which are not active
Explore which bastions are not currently active. This can be useful in identifying potential security risks or in optimizing resource usage by decommissioning inactive bastions.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_bastion_bastion
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_bastion_bastion
where
  lifecycle_state <> 'ACTIVE';
```