select title, tenant_id
from oci.oci_core_drg
where id = '{{ output.resource_id.value }}';