select title, tenant_id
from oci.oci_file_storage_snapshot
where id = '{{ output.resource_id.value }}';