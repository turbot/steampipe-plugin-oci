---
title: "Steampipe Table: oci_identity_network_source - Query OCI Identity Network Sources using SQL"
description: "Allows users to query OCI Identity Network Sources."
---

# Table: oci_identity_network_source - Query OCI Identity Network Sources using SQL

A Network Source in Oracle Cloud Infrastructure (OCI) Identity service defines a group of IP addresses that are trusted for authenticating users. It is used to limit the IP addresses that can be used to access OCI resources, adding an extra layer of security. Network Sources can be associated with groups and dynamic groups in IAM policies.

## Table Usage Guide

The `oci_identity_network_source` table provides insights into the Network Sources within OCI Identity service. As a security engineer, explore Network Source-specific details through this table, including IP address ranges, virtual source lists, and associated metadata. Utilize it to uncover information about Network Sources, such as those with specific IP ranges, the association with groups or dynamic groups, and the verification of security policies.

## Examples

### Basic info
Explore which network sources are in different lifecycle states and when they were created. This can help you manage and track your OCI identity network sources effectively.

```sql+postgres
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_identity_network_source;
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_identity_network_source;
```


### List inactive network sources
Identify network sources that are currently inactive for potential troubleshooting or resource management purposes.

```sql+postgres
select
  name,
  id,
  lifecycle_state
from
  oci_identity_network_source
where
  lifecycle_state = 'INACTIVE';
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state
from
  oci_identity_network_source
where
  lifecycle_state = 'INACTIVE';
```


### List network sources that include public IP addresses
Determine the areas in which network sources include public IP addresses. This is useful for identifying potential security vulnerabilities and ensuring proper network management.

```sql+postgres
select
  name,
  id,
  public_source_list
from
  oci_identity_network_source
where
  jsonb_array_length(public_source_list) > 0;
```

```sql+sqlite
select
  name,
  id,
  public_source_list
from
  oci_identity_network_source
where
  json_array_length(public_source_list) > 0;
```


### Get allowed VCN OCIDs and IP range pairs for each network source
Explore the allowed Virtual Cloud Network (VCN) identifiers and their corresponding IP ranges for each network source. This can help in managing and auditing network access within your cloud infrastructure.

```sql+postgres
select
  name,
  id,
  vsl ->> 'ipRanges' as ip_ranges,
  vsl ->> 'vcnId' as vcn_id
from
  oci_identity_network_source,
  jsonb_array_elements(virtual_source_list) as vsl;
```

```sql+sqlite
select
  name,
  id,
  json_extract(vsl.value, '$.ipRanges') as ip_ranges,
  json_extract(vsl.value, '$.vcnId') as vcn_id
from
  oci_identity_network_source,
  json_each(virtual_source_list) as vsl;
```