# Table: oci_dns_rrset

Domain Name System(DNS) translates domain names into IP addresses. A domain name and its matching IP address are called a DNS record.

## Examples

### Basic info

```sql
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

```sql
select
  domain,
  rtype,
  is_protected
from
  oci_dns_rrset
where
  not is_protected;
```

### List DNS records which are name server type

```sql
select
  domain,
  rtype
from
  oci_dns_rrset
where
  rtype = 'NS';
```
