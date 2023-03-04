# Table: oci_devops_project


See: https://registry.terraform.io/providers/oracle/oci/latest/docs/data-sources/devops_project 

## Examples

### Basic info

```sql
select
  display_name,
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


select   
  project.name,
  project.id, 
  project.notification_topic_id, 
  project.namespace, 
  project.description
from   
  oci_devops_project project