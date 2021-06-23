# Table: oci_dns_record

DNS or Domain Name System basically translates those domain names into IP addresses. A domain name and its matching IP address is called a “DNS record”.

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
  oci_dns_record;
```

### List DNS records which are not protected

```sql
select
  domain,
  rtype,
  is_protected
from
  oci_dns_record
where
  not is_protected;
```

### List name server type DNS records

```sql
select
  domain,
  rtype
from
  oci_dns_record
where
  rtype = 'NS';
```