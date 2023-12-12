---
title: "Steampipe Table: oci_dns_rrset - Query OCI DNS RRSet using SQL"
description: "Allows users to query DNS Resource Record Sets (RRSet) in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_dns_rrset - Query OCI DNS RRSet using SQL

A DNS Resource Record Set (RRSet) in Oracle Cloud Infrastructure (OCI) is a collection of DNS records with the same domain name, record type, and class. RRSet is used to group multiple DNS records together to provide redundancy and multiple routing paths, enhancing the performance and availability of your applications. With OCI's DNS service, you can manage RRSet to control the responses to DNS queries for your domain.

## Table Usage Guide

The `oci_dns_rrset` table provides insights into DNS Resource Record Sets (RRSet) within Oracle Cloud Infrastructure (OCI) DNS service. As a network engineer, you can explore RRSet-specific details through this table, including record type, domain name, and associated metadata. Utilize it to uncover information about RRSet, such as those with multiple routing paths, the redundancy level of your DNS records, and the verification of DNS queries for your domain.

## Examples

### Basic info
Explore which domains within your DNS records are protected, along with their type, data, and time-to-live values. This can help you assess the security and configuration of your DNS records.

```sql+postgres
select
  domain,
  rtype,
  r_data,
  ttl,
  is_protected
from
  oci_dns_rrset;
```

```sql+sqlite
select
  domain,
  rtype,
  r_data,
  ttl,
  is_protected
from
  oci_dns_rrset;
```

### List DNS records which are not protected
Explore which DNS records are not protected to bolster your system's security by identifying potential vulnerabilities. This can aid in prioritizing and implementing necessary protective measures.

```sql+postgres
select
  domain,
  rtype,
  is_protected
from
  oci_dns_rrset
where
  not is_protected;
```

```sql+sqlite
select
  domain,
  rtype,
  is_protected
from
  oci_dns_rrset
where
  not is_protected;
```

### List name server DNS records
Discover the segments that use 'NS' (Name Server) DNS records within your domain. This allows you to gain insights into your DNS configuration and understand how your domain's DNS records are managed.

```sql+postgres
select
  domain,
  rtype
from
  oci_dns_rrset
where
  rtype = 'NS';
```

```sql+sqlite
select
  domain,
  rtype
from
  oci_dns_rrset
where
  rtype = 'NS';
```