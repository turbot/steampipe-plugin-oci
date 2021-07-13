select domain, is_protected, ttl, rtype
from oci.oci_dns_rrset
where domain = '{{ output.domain.value }}' and rtype = 'NS' limit 1;