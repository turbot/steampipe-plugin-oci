---
title: "Steampipe Table: oci_network_firewall_policy - Query OCI Networking Firewall Policies using SQL"
description: "Allows users to query Firewall Policies in OCI Networking."
---

# Table: oci_network_firewall_policy - Query OCI Networking Firewall Policies using SQL

A Firewall Policy in Oracle Cloud Infrastructure (OCI) Networking is a set of rules and actions that govern the traffic flow through gateways. It provides a secure and controlled environment for network applications. Firewall Policies are essential in maintaining the security and integrity of data in OCI Networking environments.

## Table Usage Guide

The `oci_network_firewall_policy` table provides insights into Firewall Policies within OCI Networking. As a network administrator, you can explore policy-specific details through this table, including policy rules, actions, and associated metadata. Utilize it to uncover information about policies, such as those governing specific traffic, the actions associated with each policy, and the verification of policy rules.

## Examples

### Basic info
Explore the configuration and status of your network firewall policy. This information can help you assess the security rules and applications associated with your firewall, identify any mapped secrets, and determine whether the firewall is currently attached.

```sql+postgres
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

```sql+sqlite
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
Explore the recently created network firewall policies to understand their configuration and status. This is helpful to monitor the recent changes and ensure the security rules, decryption profiles, and other settings are properly configured.

```sql+postgres
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

```sql+sqlite
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
  time_created >= datetime('now', '-30 day');
```

### List network firewall policies with firewall attached
Determine the network firewall policies that have a firewall attached. This can help in identifying and managing the policies that are actively being implemented, thereby enhancing network security.

```sql+postgres
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

```sql+sqlite
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
Identify instances where network firewall policies are potentially vulnerable due to the absence of mapped secrets. This is crucial for enhancing security measures and avoiding unauthorized access.

```sql+postgres
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

```sql+sqlite
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