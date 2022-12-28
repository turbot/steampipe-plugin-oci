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