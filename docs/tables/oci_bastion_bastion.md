# Table: oci_bastion_bastion

A bastion resource. A bastion provides secured, public access to target resources in the cloud that you cannot otherwise reach from the internet. A bastion resides in a public subnet and establishes the network infrastructure needed to connect a user to a target resource in a private subnet.

## Examples

### Basic info

```sql
select
  id,
  name,
  bastion_type,
  dns_proxy_status,
  client_cidr_block_allow_list,
  max_session_ttl_in_seconds,
  max_sessions_allowed,
  private_endpoint_ip_address,
  static_jump_host_ip_address,
  phone_book_entry,
  target_vcn_id,
  target_subnet_id,
  lifecycle_state as state
from
  oci_bastion_bastion;
```

### Show Bastions that allow access from the Internet (0.0.0.0/0)

```sql
select
  id,
  name,
  bastion_type,
  client_cidr_block_allow_list,
  private_endpoint_ip_address
from
  oci_bastion_bastion
where
  (
    client_cidr_block_allow_list
  )
  ::jsonb ? '0.0.0.0/0';
```
