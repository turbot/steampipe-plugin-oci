# Table: oci_dns_tsig_key

TSIG (Transaction Signature), also referred to as Secret Key Transaction Authentication, ensures that DNS packets originate from an authorized sender by using shared secret keys and one-way hashing to add a cryptographic signature to the DNS packets. TSIG keys are used to enable DNS to authenticate updates to secondary zones. TSIG keys provide an added layer of security for IXFR and AXFR transactions.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  time_created
from
  oci_dns_tsig_key;
```

### List TSIG keys which are not active

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_dns_tsig_key
where
  lifecycle_state <> 'ACTIVE';
```
