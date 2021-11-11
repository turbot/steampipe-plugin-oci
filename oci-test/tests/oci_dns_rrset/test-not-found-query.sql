select domain, is_protected, ttl, rtype
from oci.oci_dns_rrset
where domain = 'dummy-test-{{ output.domain.value }}' and rtype = 'NS' limit 1;