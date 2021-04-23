select domain, is_protected, ttl, rtype, record_hash
from oci.oci_dns_record
where domain = '{{ output.domain.value }}';