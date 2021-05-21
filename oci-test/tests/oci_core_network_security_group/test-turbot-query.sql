select title, tenant_id, region
from oci.oci_core_network_security_group
where id = '{{ output.resource_id.value }}';