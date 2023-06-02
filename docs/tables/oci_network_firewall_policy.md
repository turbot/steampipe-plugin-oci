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

### List network firewall policies created in the last 30 days

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
  oci_network_firewall_policy
where
  time_created >= now() - interval '30' day;
```

### List network firewall policies with firewall attached

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
  oci_network_firewall_policy
where
  is_firewall_attached;
```

### List network firewall policies without mapped secrets

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
  oci_network_firewall_policy
where
  mapped_secrets is null;
```
