# Table: oci_network_firewall_policy

A network firewall policy is attached to a network firewall.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  jsonb_pretty(application_lists) as application_lists,
  jsonb_pretty(decryption_profiles) as decryption_profiles,
  jsonb_pretty(decryption_rules) as decryption_rules,
  jsonb_pretty(ip_address_lists) as ip_address_lists,
  is_firewall_attached,
  jsonb_pretty(mapped_secrets) as mapped_secrets,
  jsonb_pretty(security_rules) as security_rules,
  jsonb_pretty(url_lists) as url_lists,
  lifecycle_state as state
from
  oci_network_firewall_policy;
```

### List network firewall policies created in the last 30 days

```sql
select
  id,
  display_name,
  jsonb_pretty(application_lists) as application_lists,
  jsonb_pretty(decryption_profiles) as decryption_profiles,
  jsonb_pretty(decryption_rules) as decryption_rules,
  jsonb_pretty(ip_address_lists) as ip_address_lists,
  is_firewall_attached,
  jsonb_pretty(mapped_secrets) as mapped_secrets,
  jsonb_pretty(security_rules) as security_rules,
  jsonb_pretty(url_lists) as url_lists,
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
  jsonb_pretty(application_lists) as application_lists,
  jsonb_pretty(decryption_profiles) as decryption_profiles,
  jsonb_pretty(decryption_rules) as decryption_rules,
  jsonb_pretty(ip_address_lists) as ip_address_lists,
  is_firewall_attached,
  jsonb_pretty(mapped_secrets) as mapped_secrets,
  jsonb_pretty(security_rules) as security_rules,
  jsonb_pretty(url_lists) as url_lists,
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
  jsonb_pretty(application_lists) as application_lists,
  jsonb_pretty(decryption_profiles) as decryption_profiles,
  jsonb_pretty(decryption_rules) as decryption_rules,
  jsonb_pretty(ip_address_lists) as ip_address_lists,
  is_firewall_attached,
  jsonb_pretty(mapped_secrets) as mapped_secrets,
  jsonb_pretty(security_rules) as security_rules,
  jsonb_pretty(url_lists) as url_lists,
  lifecycle_state as state
from
  oci_network_firewall_policy
where
  mapped_secrets is null;
```
