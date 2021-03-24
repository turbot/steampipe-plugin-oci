# Table: oci_identity_policy

A policy is a document that specifies who can access which Oracle Cloud Infrastructure resources that your company has, and how. A policy simply allows a group to work in certain ways with specific types of resources in a particular compartment.

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state,
  description,
  statements
from
  oci_identity_policy;
```

### List of policies where life cycle state is not active

```sql
select
  name,
  lifecycle_state,
from
  oci_identity_policy
where lifecycle_state = 'CREATING';
```

### List of policy statements

```sql
select
  name,
  jsonb_array_elements_text(statements) as statement
from
  oci_identity_policy
```