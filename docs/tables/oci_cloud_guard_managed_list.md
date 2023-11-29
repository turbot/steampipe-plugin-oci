---
title: "Steampipe Table: oci_cloud_guard_managed_list - Query OCI Cloud Guard Managed Lists using SQL"
description: "Allows users to query OCI Cloud Guard Managed Lists."
---

# Table: oci_cloud_guard_managed_list - Query OCI Cloud Guard Managed Lists using SQL

Oracle Cloud Infrastructure (OCI) Cloud Guard is a security service that identifies security weaknesses and provides the necessary tools to investigate, track, and resolve these issues. Managed Lists in Cloud Guard are lists of resources that are managed by the service. These lists can include IP addresses, CIDR blocks, or specific users that are monitored for security threats.

## Table Usage Guide

The `oci_cloud_guard_managed_list` table provides insights into managed lists within OCI Cloud Guard. As a security analyst, you can explore details about these lists, such as the resources included, the time they were created, and their current status. Use this table to gain a comprehensive overview of the resources being monitored by Cloud Guard, helping to identify potential security threats and improve the overall security posture of your OCI environment.

## Examples

### Basic info
Explore which cloud guard managed lists have been created, when they were established, and their current lifecycle state. This can help you understand the status and timeline of your cloud guard managed lists for better resource management.

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_managed_list;
```

### List managed lists which are not active
Discover the segments that consist of managed lists in a non-active state. This can be useful to identify and assess potential areas that require attention or action in your cloud guard management.

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_managed_list
where
  lifecycle_state <> 'ACTIVE';
```