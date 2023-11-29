---
title: "Steampipe Table: oci_mysql_configuration - Query OCI MySQL Configurations using SQL"
description: "Allows users to query OCI MySQL Configurations."
---

# Table: oci_mysql_configuration - Query OCI MySQL Configurations using SQL

Oracle Cloud Infrastructure's MySQL Database Service is a fully managed database service that enables you to deploy cloud-native applications using the world’s most popular open source database. It is developed, managed, and supported by the MySQL team. This service is built on MySQL Enterprise Edition and powered by Oracle Cloud, to provide a simple, automated, integrated and enterprise ready MySQL cloud service, enabling organizations to increase business agility and reduce costs.

## Table Usage Guide

The `oci_mysql_configuration` table provides insights into MySQL configurations within Oracle Cloud Infrastructure's MySQL Database Service. As a Database Administrator, explore configuration-specific details through this table, including configuration type, shape, and associated metadata. Utilize it to uncover information about configurations, such as those related to specific MySQL versions, the details of the configuration parameters, and the verification of configuration statuses.

## Examples

### Basic info
Explore the basic details of your MySQL configurations to understand their current lifecycle state and creation time, which can be useful for managing and tracking your database resources.

```sql
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration;
```