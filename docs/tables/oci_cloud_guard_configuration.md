---
title: "Steampipe Table: oci_cloud_guard_configuration - Query OCI Cloud Guard Configurations using SQL"
description: "Allows users to query OCI Cloud Guard Configurations."
---

# Table: oci_cloud_guard_configuration - Query OCI Cloud Guard Configurations using SQL

Oracle Cloud Infrastructure's Cloud Guard is a security service that identifies security weaknesses and activities amongst your resources and assists you in rapidly acting upon them. Cloud Guard uses detectors to identify security issues and a responder to act on these issues. It monitors your entire tenancy and ensures that your resources stay secure.

## Table Usage Guide

The `oci_cloud_guard_configuration` table provides insights into Cloud Guard Configurations within Oracle Cloud Infrastructure's Cloud Guard. As a security administrator, explore configuration-specific details through this table, including status, compartment id, and associated metadata. Utilize it to uncover information about configurations, such as their current status, the compartment they belong to, and the detailed settings.

## Examples

### Basic info
Analyze the settings to understand the status and management of resources across different reporting regions in your OCI Cloud Guard configuration. This can help determine areas where resources are self-managed and where intervention may be needed.

```sql+postgres
select
  reporting_region,
  status,
  self_manage_resources
from
  oci_cloud_guard_configuration;
```

```sql+sqlite
select
  reporting_region,
  status,
  self_manage_resources
from
  oci_cloud_guard_configuration;
```