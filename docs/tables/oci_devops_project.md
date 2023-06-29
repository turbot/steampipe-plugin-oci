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
  lifecycle_state,
  time_created,
  time_updated
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
  left join oci_ons_notification_topic topic on project.notification_topic_id = topic.topic_id;
```

### List projects that are not active

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
