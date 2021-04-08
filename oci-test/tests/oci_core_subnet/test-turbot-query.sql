select title, tenant_id
from oci.oci_core_subnet
where id = '{{ output.resource_id.value }}';