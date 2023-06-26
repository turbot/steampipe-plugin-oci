# Table: oci_devops_project

A DevOps project is created to successfully build and deploy applications by using the DevOps service. A project logically groups the DevOps resources needed to implement a CI/CD workflow.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  namespace,
  state,
  time_created,
  time_updated
from
  oci_devops_project;
```

### Query with all columns for DevOps Projects

```sql
select
  name,
  id,
  lifecycle_state,
  notification_topic_id,
  namespace,
  lifecycle_details,
  description,
  time_created,
  time_updated,
  tags
from
  oci_devops_project;
```

### Join DevOps Projects with Notification Topics

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
  left join
  oci_ons_notification_topic topic
  on project.notification_topic_id = topic.topic_id
```