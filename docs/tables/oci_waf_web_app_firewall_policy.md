---
title: "Steampipe Table: oci_waf_web_app_firewall_policy - Query OCI WAF Web Application Firewall Policies using SQL"
description: "Allows users to query OCI WAF Web Application Firewall Policies."
---

# Table: oci_waf_web_app_firewall_policy - Query OCI WAF Web Application Firewall Policies using SQL

Oracle Cloud Infrastructure (OCI) Web Application Firewall (WAF) policies define the inspection and protection rules applied to HTTP traffic. A policy contains modules for request access control, request rate limiting, request protection, response access control, and response protection. Policies are attached to Web Application Firewall resources and can be shared across multiple firewalls.

## Table Usage Guide

The `oci_waf_web_app_firewall_policy` table provides insights into WAF policies within OCI. As a security engineer, you can explore policy details through this table, including configured actions, request and response control modules, and associated metadata. Use it to audit WAF rule coverage, verify that critical protection modules are enabled, and track the lifecycle of your firewall policies.

## Examples

### Basic info
Explore the WAF policies in your tenancy to understand their current lifecycle state.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_web_app_firewall_policy;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_web_app_firewall_policy;
```

### List policies not in the active state
Identify WAF policies that are not active to detect provisioning issues or resources pending deletion.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_web_app_firewall_policy
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_web_app_firewall_policy
where
  lifecycle_state <> 'ACTIVE';
```

### List policies with request access control configured
Find all WAF policies that have request access control rules defined.

```sql+postgres
select
  display_name,
  id,
  jsonb_pretty(request_access_control) as request_access_control
from
  oci_waf_web_app_firewall_policy
where
  request_access_control is not null;
```

```sql+sqlite
select
  display_name,
  id,
  request_access_control
from
  oci_waf_web_app_firewall_policy
where
  request_access_control is not null;
```

### List policies with request rate limiting enabled
Identify WAF policies that have rate limiting configured to protect against volumetric attacks.

```sql+postgres
select
  display_name,
  id,
  jsonb_pretty(request_rate_limiting) as request_rate_limiting
from
  oci_waf_web_app_firewall_policy
where
  request_rate_limiting is not null;
```

```sql+sqlite
select
  display_name,
  id,
  request_rate_limiting
from
  oci_waf_web_app_firewall_policy
where
  request_rate_limiting is not null;
```

### List policies with response protection configured
Find all WAF policies that inspect and protect outbound HTTP responses.

```sql+postgres
select
  display_name,
  id,
  jsonb_pretty(response_protection) as response_protection
from
  oci_waf_web_app_firewall_policy
where
  response_protection is not null;
```

```sql+sqlite
select
  display_name,
  id,
  response_protection
from
  oci_waf_web_app_firewall_policy
where
  response_protection is not null;
```

### List predefined actions for each policy
Review the reusable actions defined in each WAF policy for use across multiple rules.

```sql+postgres
select
  display_name,
  id,
  jsonb_pretty(actions) as actions
from
  oci_waf_web_app_firewall_policy
where
  actions is not null;
```

```sql+sqlite
select
  display_name,
  id,
  actions
from
  oci_waf_web_app_firewall_policy
where
  actions is not null;
```

### List policies by compartment
Group WAF policies by compartment to understand their distribution across your tenancy.

```sql+postgres
select
  compartment_id,
  count(*) as total,
  count(*) filter (where lifecycle_state = 'ACTIVE') as active
from
  oci_waf_web_app_firewall_policy
group by
  compartment_id;
```

```sql+sqlite
select
  compartment_id,
  count(*) as total,
  sum(case when lifecycle_state = 'ACTIVE' then 1 else 0 end) as active
from
  oci_waf_web_app_firewall_policy
group by
  compartment_id;
```
