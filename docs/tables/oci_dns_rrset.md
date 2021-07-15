# Table: oci_dns_rrset

A DNS resource record (RR) set is a collection of DNS records of the same domain and type.

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

### List name server DNS records

```sql
select
  domain,
  rtype
from
  oci_dns_rrset
where
  rtype = 'NS';
```
