# Table: oci_resourcemanager_stack

The collection of Oracle Cloud Infrastructure resources corresponding to a given Terraform configuration. Each stack resides in the compartment you specify, in a single region; however, resources on a given stack can be deployed across multiple regions. An OCID is assigned to each stack.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resourcemanager_stack;
```

### List resource manager stacks that are not active

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resourcemanager_stack
where
  lifecycle_state <> 'ACTIVE';
```