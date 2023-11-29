---
title: "Steampipe Table: oci_identity_tenancy - Query OCI Identity Tenancies using SQL"
description: "Allows users to query details about OCI Identity Tenancies."
---

# Table: oci_identity_tenancy - Query OCI Identity Tenancies using SQL

The OCI Identity Tenancy is a dedicated instance of Oracle Cloud Infrastructure resources. It is the root compartment that contains all of the organization's resources. A tenancy is provisioned in a specific Oracle Cloud Infrastructure region, but it has access to all global regions.

## Table Usage Guide

The `oci_identity_tenancy` table provides insights into the tenancies within Oracle Cloud Infrastructure Identity. As a cloud administrator, you can explore tenancy-specific details through this table, including the home region, description, and name. Utilize it to uncover information about tenancies, such as their ID, status, and time created, which can be useful for managing and auditing cloud resources.

## Examples

### Basic info
Explore the basic information about your Oracle Cloud Infrastructure (OCI) tenancy, such as its name and ID, as well as understanding its retention period and description. This is useful for getting a quick overview of your tenancy's configuration and settings.

```sql
select
  name,
  id,
  retention_period_days,
  description
from
  oci_identity_tenancy;
```

### List tenancies with a retention period less than 365 days
Explore tenancies that have a retention period of less than a year to assess compliance with data retention policies and identify any potential areas of risk.

```sql
select
  name,
  id,
  retention_period_days,
  home_region_key
from
  oci_identity_tenancy
where
  retention_period_days < 365;
```