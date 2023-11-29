---
title: "Steampipe Table: oci_devops_project - Query OCI DevOps Projects using SQL"
description: "Allows users to query DevOps Projects within Oracle Cloud Infrastructure."
---

# Table: oci_devops_project - Query OCI DevOps Projects using SQL

A DevOps Project in Oracle Cloud Infrastructure (OCI) is a managed environment that provides a set of resources for developing, testing, and delivering applications in the cloud. It offers a unified interface for managing both infrastructure as code and application code, helping developers and operations teams work together more effectively. DevOps Projects in OCI enable continuous integration and delivery, infrastructure automation, and other DevOps practices.

## Table Usage Guide

The `oci_devops_project` table provides insights into DevOps Projects within Oracle Cloud Infrastructure. As a DevOps engineer, you can explore project-specific details through this table, including project details, associated resources, and other metadata. Utilize it to uncover information about projects, such as those with specific configurations, the relationships between resources within a project, and the status of different elements of the project.

## Examples

### Basic info
Explore the status and timing of your DevOps projects. This query allows you to monitor project progress and updates, providing insight into the lifecycle and timelines of your projects.

```sql
select
  name,
  id,
  description,
  namespace,
  lifecycle_state,
  time_created,
  time_updated
from
  oci_devops_project;
```

### Join DevOps Projects with Notification Topics
Explore the connection between DevOps projects and their corresponding notification topics. This can be crucial in managing project notifications and understanding the overall project configuration.

```sql
select
  project.name,
  project.id,
  project.notification_topic_id,
  topic.name notification_topic_name,
  project.namespace,
  project.description
from
  oci_devops_project project
  left join oci_ons_notification_topic topic on project.notification_topic_id = topic.topic_id;
```

### List projects that are not active
Discover the segments that consist of inactive projects within your organisation. This can be useful to identify and manage projects that are no longer in use, thereby optimizing resources and improving operational efficiency.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created,
  namespace
from
  oci_devops_project
where
  lifecycle_state <> 'ACTIVE';
```

### List projects that are created in the last 30 days
Explore which projects have been initiated in the recent 30 days. This is useful for tracking recent development activities and ensuring necessary follow-ups.

```sql
select
  name,
  id,
  time_created,
  namespace
from
  oci_devops_project
where
  time_created >= now() - interval '30' day;
```