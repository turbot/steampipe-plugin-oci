---
title: "Steampipe Table: oci_identity_availability_domain - Query OCI Identity Availability Domains using SQL"
description: "Allows users to query OCI Identity Availability Domains."
---

# Table: oci_identity_availability_domain - Query OCI Identity Availability Domains using SQL

Oracle Cloud Infrastructure's Identity and Access Management (IAM) service lets you control who has access to your cloud resources. Availability Domains are isolated from each other, fault-tolerant, and very unlikely to fail simultaneously. They provide the resources to run your applications and databases in a high-availability configuration.

## Table Usage Guide

The `oci_identity_availability_domain` table provides insights into Availability Domains within OCI Identity and Access Management (IAM). As a cloud engineer, you can explore domain-specific details through this table, including their names, compartment IDs, and associated metadata. Use it to uncover information about domains, such as their status, the resources they hold, and their relationships with other domains.

## Examples

### Basic info
Explore the names and IDs of available domains within the Oracle Cloud Infrastructure (OCI) to manage and organize resources more efficiently. This can be particularly useful for system administrators seeking to streamline their resource allocation and tracking processes.

```sql+postgres
select
  name,
  id
from
  oci_identity_availability_domain;
```

```sql+sqlite
select
  name,
  id
from
  oci_identity_availability_domain;
```