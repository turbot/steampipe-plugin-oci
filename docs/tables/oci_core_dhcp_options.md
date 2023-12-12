---
title: "Steampipe Table: oci_core_dhcp_options - Query OCI Core DHCP Options using SQL"
description: "Allows users to query OCI Core DHCP Options."
---

# Table: oci_core_dhcp_options - Query OCI Core DHCP Options using SQL

OCI Core DHCP Options is a resource within Oracle Cloud Infrastructure that allows you to manage options for the Dynamic Host Configuration Protocol (DHCP). DHCP options determine how DHCP-enabled devices operate in a VCN (Virtual Cloud Network). It includes options such as domain name servers, search domains, and more, providing a centralized way to manage and configure these options for your VCN.

## Table Usage Guide

The `oci_core_dhcp_options` table provides insights into DHCP Options within Oracle Cloud Infrastructure Core Services. As a network administrator, explore option-specific details through this table, including the associated VCN, domain name servers, search domains, and more. Utilize it to manage and configure your network settings, ensuring optimal operation of DHCP-enabled devices within your VCN.

## Examples

### Basic info
Explore the state and creation time of your Oracle Cloud Infrastructure's DHCP options to understand their lifecycle and location. This can help you assess their configuration and manage your resources more efficiently.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  region
from
  oci_core_dhcp_options;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  region
from
  oci_core_dhcp_options;
```


### Get configuration info for each DHCP option
Explore the configuration details of each DHCP option to understand the server type and custom DNS servers. This can be particularly useful for network administrators who want to manage and optimize their network configurations.

```sql+postgres
select
  id,
  display_name,
  jsonb_array_elements_text(o -> 'searchDomainNames') as search_domain_names,
  jsonb_array_elements_text(o -> 'customDnsServers') as custom_dns_servers,
  o ->> 'serverType' as server_type,
  o ->> 'type' as type
from
  oci_core_dhcp_options,
  jsonb_array_elements(options) as o;
```

```sql+sqlite
select
  id,
  display_name,
  json_extract(o.value, '$.searchDomainNames') as search_domain_names,
  json_extract(o.value, '$.customDnsServers') as custom_dns_servers,
  json_extract(o.value, '$.serverType') as server_type,
  json_extract(o.value, '$.type') as type
from
  oci_core_dhcp_options,
  json_each(options) as o;
```


### Count the number of DHCP options by VCN
Identify the quantity of DHCP options for each Virtual Cloud Network (VCN) to understand the network configuration and manage resources effectively. This can aid in optimizing network performance and troubleshooting network issues.

```sql+postgres
select
  vcn_id,
  count(*) dhcp_option_count
from
  oci_core_dhcp_options
group by
  vcn_id;
```

```sql+sqlite
select
  vcn_id,
  count(*) dhcp_option_count
from
  oci_core_dhcp_options
group by
  vcn_id;
```