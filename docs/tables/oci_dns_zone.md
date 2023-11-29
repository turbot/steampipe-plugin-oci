---
title: "Steampipe Table: oci_dns_zone - Query OCI DNS Zones using SQL"
description: "Allows users to query OCI DNS Zones."
---

# Table: oci_dns_zone - Query OCI DNS Zones using SQL

Oracle Cloud Infrastructure's DNS service allows you to manage and control your DNS zones. It provides a global, scalable, high-performance DNS solution that can direct your end users to your internet applications by translating domain names into IP addresses. This service supports various types of DNS records, including A, AAAA, CNAME, MX, PTR, and TXT.

## Table Usage Guide

The `oci_dns_zone` table provides insights into DNS zones within Oracle Cloud Infrastructure's DNS service. As a network administrator, explore zone-specific details through this table, including the zone's name, state, and associated records. Utilize it to uncover information about zones, such as those with specific types of records, the time to live (TTL) settings for each record, and the verification of zone configurations.

## Examples

### Basic info
Discover the segments of your OCI DNS zones that are currently active or inactive. This can help you understand when each zone was created, allowing you to manage and optimize your resources effectively.

```sql
select
  id,
  name,
  lifecycle_state,
  time_created
from
  oci_dns_zone;
```

### List DNS zones which are not active
Explore DNS zones that are not currently active, to identify potential issues or areas for clean-up. This can help in maintaining an efficient and clean DNS system.

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_dns_zone
where
  lifecycle_state in ('CREATING','DELETED','DELETING','FAILED');
```