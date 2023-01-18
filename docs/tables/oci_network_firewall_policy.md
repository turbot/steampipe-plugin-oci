# Table: oci_network_firewall_policy

A network firewall policy is attached to a network firewall.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  application_lists,
  decryption_profiles,
  decryption_rules,
  ip_address_lists,
  is_firewall_attached,
  mapped_secrets,
  security_rules,
  url_lists,
  lifecycle_state as state 
from
  oci_network_firewall_policy;
```
