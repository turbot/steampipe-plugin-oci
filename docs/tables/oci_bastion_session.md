# Table: oci_bastion_session

A bastion session resource. A bastion session lets authorized users connect to a target resource using a Secure Shell (SSH) for a predetermined amount of time.

## Examples

### Basic info

```sql
select
  s.id,
  s.bastion_id,
  s.display_name,
  s.bastion_name,
  s.target_resource_details,
  s.key_details,
  s.session_ttl_in_seconds,
  s.bastion_user_name,
  s.ssh_metadata,
  s.key_type,
  s.lifecycle_state as state 
from
  oci_bastion_session s 
  INNER JOIN
    oci_bastion_bastion b 
    ON b.id = s.bastion_id
```

### Show port forwarding bastion sessions

```sql
select
  id,
  bastion_id,
  display_name,
  bastion_name,
  target_resource_details,
  key_details,
  session_ttl_in_seconds,
  bastion_user_name,
  ssh_metadata,
  key_type,
  lifecycle_state as state 
from
  oci_bastion_session 
where
  bastion_id = 'ocid' 
  and target_resource_details -> 'sessionType' = '"MANAGED_SSH"'
```
