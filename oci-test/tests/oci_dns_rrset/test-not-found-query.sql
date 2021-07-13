select domain, is_protected, ttl
from oci.oci_dns_rrset
where domain = '{{ output.domain.value }}aa';