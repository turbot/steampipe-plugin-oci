select title, tenant_id, region
from oci.oci_file_storage_mount_target
where id = '{{ output.resource_id.value }}';