---
title: "Steampipe Table: oci_waf_web_app_firewall - Query OCI WAF Web Application Firewalls using SQL"
description: "Allows users to query OCI WAF Web Application Firewalls."
---

# Table: oci_waf_web_app_firewall - Query OCI WAF Web Application Firewalls using SQL

Oracle Cloud Infrastructure (OCI) Web Application Firewall (WAF) protects HTTP-based applications from web attacks. A Web Application Firewall resource attaches a WAF policy to a backend resource. Currently, OCI WAF supports the `LOAD_BALANCER` backend type, linking a WAF policy to an OCI Load Balancer to inspect and filter inbound and outbound HTTP traffic.

## Table Usage Guide

The `oci_waf_web_app_firewall` table provides insights into WAF instances within OCI. As a security engineer, you can explore firewall-specific details through this table, including the attached policy, backend type, and load balancer association. Use it to audit which load balancers are protected by WAF, verify policy assignments, and track the lifecycle of your firewall resources.

## Examples

### Basic info
Explore the web application firewalls in your tenancy to understand which load balancers are protected and which policies are applied.

```sql+postgres
select
  display_name,
  id,
  backend_type,
  web_app_firewall_policy_id,
  lifecycle_state
from
  oci_waf_web_app_firewall;
```

```sql+sqlite
select
  display_name,
  id,
  backend_type,
  web_app_firewall_policy_id,
  lifecycle_state
from
  oci_waf_web_app_firewall;
```

### List firewalls not in the active state
Identify web application firewalls that are not active to detect provisioning issues or resources pending deletion.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_web_app_firewall
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_web_app_firewall
where
  lifecycle_state <> 'ACTIVE';
```

### List firewalls with their associated load balancer
Identify which load balancers have WAF protection enabled by joining firewall resources to their backend load balancer IDs.

```sql+postgres
select
  display_name,
  id,
  load_balancer_id,
  web_app_firewall_policy_id
from
  oci_waf_web_app_firewall
where
  backend_type = 'LOAD_BALANCER';
```

```sql+sqlite
select
  display_name,
  id,
  load_balancer_id,
  web_app_firewall_policy_id
from
  oci_waf_web_app_firewall
where
  backend_type = 'LOAD_BALANCER';
```

### List firewalls using a specific WAF policy
Find all firewalls that share a particular WAF policy to understand the blast radius of policy changes.

```sql+postgres
select
  display_name,
  id,
  backend_type,
  lifecycle_state
from
  oci_waf_web_app_firewall
where
  web_app_firewall_policy_id = 'ocid1.webappfirewallpolicy.oc1..example';
```

```sql+sqlite
select
  display_name,
  id,
  backend_type,
  lifecycle_state
from
  oci_waf_web_app_firewall
where
  web_app_firewall_policy_id = 'ocid1.webappfirewallpolicy.oc1..example';
```

### List firewalls by compartment
Group firewalls by compartment to understand their distribution across your tenancy.

```sql+postgres
select
  compartment_id,
  count(*) as total,
  count(*) filter (where lifecycle_state = 'ACTIVE') as active
from
  oci_waf_web_app_firewall
group by
  compartment_id;
```

```sql+sqlite
select
  compartment_id,
  count(*) as total,
  sum(case when lifecycle_state = 'ACTIVE' then 1 else 0 end) as active
from
  oci_waf_web_app_firewall
group by
  compartment_id;
```
