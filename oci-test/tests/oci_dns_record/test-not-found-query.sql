select domain, is_protected, ttl
from oci.oci_dns_record
where domain = '{{ output.domain.value }}aa';