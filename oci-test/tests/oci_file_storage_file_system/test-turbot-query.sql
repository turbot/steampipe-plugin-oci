select title, tenant_id, region
from oci.oci_file_storage_file_system
where id = '{{ output.resource_id.value }}';