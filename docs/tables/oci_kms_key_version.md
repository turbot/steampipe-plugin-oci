---
title: "Steampipe Table: oci_kms_key_version - Query OCI Key Management Service Key Versions using SQL"
description: "Allows users to query Key Versions within the Oracle Cloud Infrastructure Key Management Service."
---

# Table: oci_kms_key_version - Query OCI Key Management Service Key Versions using SQL

The Oracle Cloud Infrastructure Key Management service is a cloud solution that lets you manage the cryptographic keys used to protect your data. Key versions are immutable, meaning once they are created, they can't be changed. They are used to encrypt and decrypt data, and each key version is associated with a master key.

## Table Usage Guide

The `oci_kms_key_version` table provides insights into Key Versions within the Oracle Cloud Infrastructure Key Management Service. As a security engineer, explore key version-specific details through this table, including their lifecycle states, creation times, and associated master keys. Utilize it to uncover information about key versions, such as those with specific cryptographic endpoints, the association between key versions and master keys, and the status of each key version.

## Examples

### Basic info
Explore the lifecycle status and creation time of specific keys in your Oracle Cloud Infrastructure Key Management service. This can be useful in managing and tracking your encryption keys, ensuring they are in the desired state and were created at the expected time.

```sql
select
  v.id as key_version_id,
  k.name as key_name,
  v.lifecycle_state,
  v.time_created as time_created
from
  oci_kms_key k,
  oci_kms_key_version v
where
  v.key_id = k.id
  and v.management_endpoint = k.management_endpoint
  and v.region = k.region;
```

### Get latest key version for all active keys
Identify the most recent versions of all active keys in your system. This could be useful for auditing purposes or to ensure that you're always using the most up-to-date keys for security purposes.

```sql
with oci_kms as (
  select
    k.name,
    k.lifecycle_state,
    max(v.time_created) as latest_key_version_created,
    k.region,
    k.compartment_id
  from
    oci_kms_key k, oci_kms_key_version v
  where 
    v.key_id = k.id
    and v.management_endpoint = k.management_endpoint
    and v.region = k.region
    and k.lifecycle_state = 'ENABLED'
  group by
    k.name,k.lifecycle_state, k.region,k.compartment_id
)
select 
  k.name,
  k.lifecycle_state,
  latest_key_version_created,
  k.region,
  coalesce(c.name, 'root') as compartment
from
  oci_kms k left join oci_identity_compartment c on c.id = k.compartment_id;
```