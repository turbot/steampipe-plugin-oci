---
title: "Steampipe Table: oci_certificates_management_ca_bundle - Query OCI Certificate Management CA Bundles using SQL"
description: "Allows users to query Certificate Management CA Bundles in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_certificates_management_ca_bundle - Query OCI Certificate Management CA Bundles using SQL

Oracle Cloud Infrastructure (OCI) Certificate Management is a feature that allows you to manage and deploy SSL/TLS certificates. It provides a centralized way to manage certificates, ensuring secure and encrypted communication between the client and the server. OCI Certificate Management helps in maintaining the trust and integrity of your applications and infrastructure.

## Table Usage Guide

The `oci_certificates_management_ca_bundle` table provides insights into Certificate Management CA Bundles within OCI. As a Security Engineer, explore bundle-specific details through this table, including certificate data, validity, and associated metadata. Utilize it to uncover information about certificate bundles, such as those with expired certificates, the trust relationships between certificates, and the verification of trust policies.

## Examples

### Basic info
Explore which certificates are active or inactive by analyzing their lifecycle state and details. This can help you understand the status and creation time of certificates for better management and timely renewal.

```sql
select
  id,
  name,
  lifecycle_state,
  lifecycle_details,
  description,
  time_created
from
  oci_certificates_management_ca_bundle;
```

### List bundles created between a specific time
Discover the segments that were created within a specific time frame. This query is useful in tracking the lifecycle of your resources, allowing you to manage and optimize their utilization effectively.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_certificates_management_ca_bundle
where
  time_created between '2023-05-01' and '2023-07-01';
```

### List failed bundles
Determine the areas in which certificate management bundles have failed, allowing you to identify and address issues in your OCI environment.

```sql
select
  name,
  id,
  lifecycle_state,
  lifecycle_details
from
  oci_certificates_management_ca_bundle
where
  lifecycle_state = 'FAILED';
```