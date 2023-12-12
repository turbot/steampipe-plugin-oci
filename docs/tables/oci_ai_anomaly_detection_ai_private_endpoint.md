---
title: "Steampipe Table: oci_ai_anomaly_detection_ai_private_endpoint - Query OCI AI Anomaly Detection AI Private Endpoints using SQL"
description: "Allows users to query OCI AI Anomaly Detection AI Private Endpoints."
---

# Table: oci_ai_anomaly_detection_ai_private_endpoint - Query OCI AI Anomaly Detection AI Private Endpoints using SQL

AI Private Endpoint is a feature in Oracle Cloud Infrastructure's AI Anomaly Detection service. It provides a private, secure, and direct connection between your virtual cloud network (VCN) and the AI services, without exposing the data to the public internet. This feature enhances the security of your data by keeping all traffic within the Oracle network.

## Table Usage Guide

The `oci_ai_anomaly_detection_ai_private_endpoint` table provides insights into AI Private Endpoints within Oracle Cloud Infrastructure's AI Anomaly Detection service. As a security analyst or a data engineer, explore endpoint-specific details through this table, including the associated network details, the status of the endpoint, and other metadata. Utilize it to uncover information about endpoints, such as their current status, the network they are associated with, and the time they were created.

## Examples

### Basic info
Explore which anomaly detection endpoints are currently active in your Oracle Cloud Infrastructure. This can help you understand the state of your AI services and manage resources effectively.

```sql+postgres
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint;
```

```sql+sqlite
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint;
```

### List the AI private endpoints created in the last 30 days
Explore the recent additions to your AI private endpoints over the past month. This can help in tracking the growth and changes in your AI infrastructure, ensuring you stay updated on the evolution of your system.

```sql+postgres
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  time_created >= datetime('now', '-30 day');
```

### List the AI private endpoints that have attached data assets
Discover the segments that have AI private endpoints associated with data assets. This is useful in understanding which endpoints are actively being used for data processing in anomaly detection, providing insights for resource allocation and security monitoring.

```sql+postgres
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  attached_data_assets is not null;
```

```sql+sqlite
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  attached_data_assets is not null;
```

### List AI private endpoints which are not active
Discover the segments that consist of inactive AI private endpoints. This can be beneficial in identifying endpoints that may need troubleshooting or updating.

```sql+postgres
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  lifecycle_state <> 'ACTIVE';
```

### List DNS zones associated to the AI private endpoints
Explore the association of DNS zones with AI private endpoints in your network. This can help in identifying potential anomalies and understanding the lifecycle state of these endpoints.

```sql+postgres
select
  id,
  subnet_id,
  jsonb_pretty(dns_zones) as dns_zones,
  display_name,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  dns_zones is not null;
```

```sql+sqlite
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  dns_zones is not null;
```