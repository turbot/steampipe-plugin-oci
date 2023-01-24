select title, tenant_id
from oci.oci_core_volume_group
where id = '{{ output.resource_id.value }}';