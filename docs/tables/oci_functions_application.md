---
title: "Steampipe Table: oci_functions_application - Query OCI Functions Applications using SQL"
description: "Allows users to query data related to Applications within the Oracle Cloud Infrastructure (OCI) Functions service."
---

# Table: oci_functions_application - Query OCI Functions Applications using SQL

Oracle Cloud Infrastructure Functions is a fully managed, multi-tenant, highly scalable, on-demand, Functions-as-a-Service platform. It allows developers to write and deploy code with functions that are triggered by events or direct invocations. Applications in OCI Functions represent a logical grouping of functions, and the `oci_functions_application` table can be used to query these Applications.

## Table Usage Guide

The `oci_functions_application` table provides insights into Applications within the Oracle Cloud Infrastructure Functions service. As a developer or DevOps engineer, explore Application-specific details through this table, including the configuration and state of each Application. Utilize it to uncover information about the Applications, such as their associated functions, the network configuration, and the compartment they belong to.

## Examples

### Basic info
Discover the status and identifiers of your Oracle Cloud Infrastructure's functions applications. This allows for a quick overview and understanding of your application's lifecycle state and associated subnet IDs.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  lifecycle_state,
  subnet_ids
from
  oci_functions_application;
```


### List applications not in the active state
Determine the areas in which applications are not currently active. This can help in identifying applications that may need attention or troubleshooting, thereby ensuring smooth operations.

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_functions_application
where
  lifecycle_state <> 'ACTIVE';
```


### Get configuration details for each application
Analyze the settings to understand each application's configuration within Oracle Cloud Infrastructure's Functions service. This allows for a comprehensive review of application settings, aiding in optimization and troubleshooting efforts.

```sql
select
  display_name,
  id,
  config
from
  oci_functions_application;
```