select title, tenant_id
from oci.oci_database_software_image
where id = '{{ output.resource_id.value }}';