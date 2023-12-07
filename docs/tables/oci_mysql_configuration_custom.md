---
title: "Steampipe Table: oci_mysql_configuration_custom - Query OCI MySQL Custom Configurations using SQL"
description: "Allows users to query custom configurations of the OCI MySQL Database Service."
---

# Table: oci_mysql_configuration_custom - Query OCI MySQL Custom Configurations using SQL

The Oracle Cloud Infrastructure (OCI) MySQL Database Service is a fully managed database service that enables organizations to deploy cloud-native applications using the world's most popular open source database. It is developed, managed, and maintained by the same team that develops MySQL, ensuring up-to-date features, compatibility, and enhancements. A custom configuration in OCI MySQL Database Service allows users to customize various MySQL system variables to suit their specific requirements.

## Table Usage Guide

The `oci_mysql_configuration_custom` table provides insights into custom configurations within OCI MySQL Database Service. As a database administrator or developer, explore configuration-specific details through this table, including the configuration's ID, name, description, shape name, and various MySQL variables. Utilize it to uncover information about specific custom configurations, such as their variable settings, the MySQL version they're compatible with, and their associated compartment ID.

## Examples

### Basic info
Uncover the details of custom MySQL configurations in your Oracle Cloud Infrastructure, including their display names, IDs, descriptions, lifecycle states, and creation times. This can help you manage and monitor your custom configurations more effectively.

```sql+postgres
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration_custom;
```

```sql+sqlite
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration_custom;
```

### List deleted configurations
Explore which custom MySQL configurations have been deleted in your Oracle Cloud Infrastructure. This is useful for tracking changes and maintaining the security and efficiency of your database systems.

```sql+postgres
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration_custom 
where
  lifecycle_state = 'DELETED';
```

```sql+sqlite
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration_custom 
where
  lifecycle_state = 'DELETED';
```