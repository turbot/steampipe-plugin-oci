---
title: "Steampipe Table: oci_core_service_gateway - Query OCI Core Services using SQL"
description: "Allows users to query Service Gateways in OCI Core Services."
---

# Table: oci_core_service_gateway - Query OCI Core Services using SQL

A Service Gateway in Oracle Cloud Infrastructure (OCI) Core Services enables your on-premises network to access Oracle services without exposing the data to the public internet. It provides private access from your virtual cloud network (VCN) to services in the Oracle Services Network. The traffic from your VCN to the Oracle service travels over the Oracle network fabric and never traverses the internet.

## Table Usage Guide

The `oci_core_service_gateway` table provides insights into the Service Gateways within OCI Core Services. As a network administrator, you can use this table to explore details about each Service Gateway, including its associated VCN, lifecycle state, and the Oracle services enabled for the gateway. This can be beneficial for monitoring the status and configuration of your private network access to Oracle services.

## Examples

### Basic info
Explore which service gateways in your Oracle Cloud Infrastructure are active and when they were created. This can help you manage your resources and understand your usage patterns across different regions.

```sql
select
  display_name,
  id,
  time_created,
  vcn_id,
  lifecycle_state,
  region
from
  oci_core_service_gateway;
```

### List service gateways that use route tables
Explore which service gateways are utilizing route tables. This can be beneficial in identifying potential areas of network configuration that may require optimization or troubleshooting.

```sql
select
  display_name,
  id,
  route_table_id
from
  oci_core_service_gateway
where
  route_table_id is not null;
```

### Get enabled services for each service gateway
Explore which services are currently enabled for each service gateway. This can be useful in managing and optimizing network traffic by identifying active services.

```sql
select
  display_name,
  id,
  s ->> 'serviceId' as service_id,
  s ->> 'serviceName' as service_name
from
  oci_core_service_gateway,
  jsonb_array_elements(services) as s;
```

### List service gateways that block traffic
Discover the segments that are obstructing traffic flow within the service gateways. This is useful in identifying potential bottlenecks or areas of concern within your network infrastructure.

```sql
select
  display_name,
  id,
  block_traffic
from
  oci_core_service_gateway
where
  block_traffic;
```