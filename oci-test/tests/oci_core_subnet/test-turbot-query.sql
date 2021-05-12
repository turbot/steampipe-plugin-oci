select title, tenant_id, region
from oci.oci_core_subnet
where id = '{{ output.resource_id.value }}';