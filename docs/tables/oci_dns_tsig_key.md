---
title: "Steampipe Table: oci_dns_tsig_key - Query OCI DNS TSIG Keys using SQL"
description: "Allows users to query DNS Transaction Signature (TSIG) Keys in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_dns_tsig_key - Query OCI DNS TSIG Keys using SQL

DNS Transaction Signature (TSIG) is a resource in Oracle Cloud Infrastructure (OCI) that provides enhanced security for DNS operations. TSIG keys are used to authenticate DNS messages and prevent unauthorized DNS operations. They are crucial for maintaining the integrity and security of DNS data within OCI.

## Table Usage Guide

The `oci_dns_tsig_key` table offers insights into DNS Transaction Signature (TSIG) Keys within OCI DNS. As a security analyst, you can explore key-specific details through this table, including algorithm types, secret keys, and associated metadata. Utilize it to uncover information about TSIG keys, such as those with specific algorithms, the lifecycle state of keys, and the time of their creation.

## Examples

### Basic info
Assess the elements within your Oracle Cloud Infrastructure by pinpointing specific locations where certain security keys were created and their current lifecycle state. This aids in understanding the status and age of your security measures.

```sql
select
  id,
  name,
  lifecycle_state,
  time_created
from
  oci_dns_tsig_key;
```

### List TSIG keys which are not active
Explore TSIG keys that are not currently active. This can be useful in identifying keys that may be expired or unused, aiding in the maintenance and security of your DNS system.

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_dns_tsig_key
where
  lifecycle_state <> 'ACTIVE';
```