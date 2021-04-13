select title, tenant_id
from oci.oci_core_volume_backup
where id = '{{ output.resource_id.value }}';