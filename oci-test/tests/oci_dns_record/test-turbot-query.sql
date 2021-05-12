select title, tenant_id
from oci.oci_dns_record
where domain = '{{ output.domain.value }}' limit 1;