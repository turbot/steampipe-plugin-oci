select title, tenant_id, region
from oci.oci_functions_application
where id = '{{ output.resource_id.value }}';