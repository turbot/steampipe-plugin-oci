select title, tenant_id
from oci.oci_core_image
where id = '{{ output.resource_id.value }}';