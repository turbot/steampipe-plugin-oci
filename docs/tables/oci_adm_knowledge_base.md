---
title: "Steampipe Table: oci_adm_knowledge_base - Query OCI ADM Knowledge Bases using SQL"
description: "Allows users to query OCI ADM Knowledge Bases."
---

# Table: oci_adm_knowledge_base - Query OCI ADM Knowledge Bases using SQL

Oracle Cloud Infrastructure (OCI) Audit and Data Management (ADM) Knowledge Base is a cloud-based service that provides a centralized repository for managing and accessing information related to your OCI resources. This service is designed to help users manage, share, and access knowledge across their organization, including information about OCI resources, incidents, and solutions. It enables users to efficiently retrieve information, reduce resolution time, and improve operational efficiency.

## Table Usage Guide

The `oci_adm_knowledge_base` table provides insights into the Knowledge Bases within Oracle Cloud Infrastructure's Audit and Data Management (ADM). As a cloud administrator or a DevOps engineer, you can explore detailed information about your organization's knowledge bases, including the associated incidents, solutions, and other related resources. Utilize this table to maintain an efficient knowledge management system, streamline incident resolution, and enhance operational efficiency within your organization.

## Examples

### Basic info
Explore which knowledge base entries are active or inactive within your Oracle Cloud Infrastructure. This can help you manage and organize your resources effectively.

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base;
```

### List knowledge bases which are not active
Discover the segments that consist of inactive knowledge bases. This can be particularly useful when carrying out maintenance or auditing tasks, as it helps identify areas that may require attention or updates.

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base
where
  lifecycle_state <> 'ACTIVE';
```

### List knowledge bases created in last 30 days
Discover the knowledge bases that have been created within the past 30 days to maintain an up-to-date understanding of your data resources and ensure you're tracking recent developments.

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base
where
  time_created >= now() - interval '30' day;
```

### List knowledge bases that have not been updated for more than 90 days
Discover the segments that have not been frequently updated, specifically knowledge bases that have been dormant for more than 90 days. This can be useful for identifying potentially outdated or unused resources that may require attention or cleanup.

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base
where
  time_updated < now() - interval '90' day;
```