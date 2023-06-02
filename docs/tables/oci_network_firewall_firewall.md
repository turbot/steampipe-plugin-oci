# Table: oci_network_firewall_firewall

A network firewall is a security resource that exists in a subnet of your choice and controls incoming and outgoing network traffic based on a set of security rules. Each firewall is associated with a policy. Traffic is routed to and from the firewall from resources such as internet gateways and dynamic routing gateways (DRGs).

## Examples

### Basic info

```sql
select
  id,
  display_name,
  availability_domain,
  ipv4_address,
  ipv6_address,
  network_firewall_policy_id,
  network_security_group_ids,
  subnet_id,
  lifecycle_state as state
from
  oci_network_firewall_firewall;
```

### List network firewalls created in the last 30 days

```sql
select
  id,
  display_name,
  availability_domain,
  ipv4_address,
  ipv6_address,
  network_firewall_policy_id,
  network_security_group_ids,
  subnet_id,
  lifecycle_state as state
from
  oci_network_firewall_firewall
where
  time_created >= now() - interval '30' day;
```

### List network firewalls having IPv6 address

```sql
select
  id,
  display_name,
  availability_domain,
  ipv4_address,
  ipv6_address,
  network_firewall_policy_id,
  network_security_group_ids,
  subnet_id,
  lifecycle_state as state
from
  oci_network_firewall_firewall
where
  ipv6_address is not null;
```

### Describe the network firewall policy associated to the network firewall

```sql
select
  f.display_name as firewall_name,
  f.id as firewall_id,
  p.display_name as policy_display_name,
  p.id as policy_id,
  p.lifecycle_details as policy_lifecycle 
from
  oci_network_firewall_firewall as f 
  left join
    oci_network_firewall_policy as p 
    on f.network_firewall_policy_id = p.id;
```
